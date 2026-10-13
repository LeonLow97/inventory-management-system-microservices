#!/usr/bin/env bash

awscli_load_config() {
  local root_dir
  root_dir="$(cd "${1:?AWS CLI root is required}" && pwd)"

  if [[ ! -f "$root_dir/.env.template" ]]; then
    echo "Missing $root_dir/.env.template." >&2
    return 1
  fi

  set -a
  # Load defaults first, then local overrides.
  source "$root_dir/.env.template"
  if [[ -f "$root_dir/.env" ]]; then
    source "$root_dir/.env"
  fi
  set +a

  VARIABLES_FILE="${VARIABLES_FILE:-variables.env}"
  if [[ "$VARIABLES_FILE" != /* ]]; then
    VARIABLES_FILE="$root_dir/$VARIABLES_FILE"
  fi
  touch "$VARIABLES_FILE"
  chmod 600 "$VARIABLES_FILE"

  if [[ -z "${REGION:-}" ]]; then
    echo "REGION is missing from $root_dir/.env." >&2
    return 1
  fi
  # Pin commands without an explicit --region to the project configuration.
  export AWS_REGION="$REGION"
  export AWS_DEFAULT_REGION="$REGION"
}

awscli_run_allowing_error_code() {
  local expected_error_code="${1:?Expected AWS error code is required}"
  shift

  local output
  if output="$("$@" 2>&1)"; then
    [[ -z "$output" ]] || printf '%s\n' "$output"
    return 0
  fi

  if [[ "$output" == *"($expected_error_code)"* ]]; then
    return 0
  fi

  printf '%s\n' "$output" >&2
  return 1
}

awscli_load_state() {
  set -a
  source "${VARIABLES_FILE:?Call awscli_load_config first}"
  set +a
}

awscli_rewrite_state() {
  local operation="$1"
  shift

  local entry key value escaped_value key_list="" state_file temporary_file
  local keys=()
  for entry in "$@"; do
    key="${entry%%=*}"
    keys+=("$key")
  done
  key_list="$(IFS=,; printf '%s' "${keys[*]}")"

  state_file="${VARIABLES_FILE:?Call awscli_load_config first}"
  temporary_file="$(mktemp "${state_file}.tmp.XXXXXX")"
  if ! awk -F= -v key_list="$key_list" '
    BEGIN {
      count = split(key_list, key_names, ",")
      for (i = 1; i <= count; i++) removed[key_names[i]] = 1
    }
    !($1 in removed)
  ' "$state_file" > "$temporary_file"; then
    rm -f "$temporary_file"
    return 1
  fi

  if [[ "$operation" == "write" ]]; then
    for entry in "$@"; do
      key="${entry%%=*}"
      value="${entry#*=}"
      if [[ ! "$key" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]]; then
        rm -f "$temporary_file"
        echo "Invalid state key '$key'." >&2
        return 1
      fi
      printf -v escaped_value '%q' "$value"
      printf '%s=%s\n' "$key" "$escaped_value" >> "$temporary_file"
    done
  fi
  mv "$temporary_file" "$state_file"
  chmod 600 "$state_file"
}

awscli_write_state() {
  awscli_rewrite_state write "$@"
}

awscli_remove_state() {
  awscli_rewrite_state remove "$@"
}

awscli_wait_for_ssm_instance() {
  local instance_id="${1:?Instance ID is required}"
  local attempts="${2:-60}"
  local attempt status

  for ((attempt = 1; attempt <= attempts; attempt++)); do
    status="$(aws ssm describe-instance-information \
      --filters "Key=InstanceIds,Values=$instance_id" \
      --query 'InstanceInformationList[0].PingStatus' \
      --output text --region "$REGION" 2>/dev/null || true)"
    if [[ "$status" == "Online" ]]; then
      return 0
    fi
    sleep 5
  done

  echo "SSM did not report instance $instance_id online after $((attempts * 5)) seconds." >&2
  return 1
}

awscli_ensure_instance_running() {
  local instance_id="${1:?Instance ID is required}"
  local state

  state="$(aws ec2 describe-instances \
    --instance-ids "$instance_id" \
    --query 'Reservations[0].Instances[0].State.Name' \
    --output text --region "$REGION")"

  case "$state" in
    stopped)
      aws ec2 start-instances --instance-ids "$instance_id" --region "$REGION" >/dev/null
      ;;
    stopping)
      aws ec2 wait instance-stopped --instance-ids "$instance_id" --region "$REGION"
      aws ec2 start-instances --instance-ids "$instance_id" --region "$REGION" >/dev/null
      ;;
    pending|running)
      ;;
    *)
      echo "Cannot deploy to instance $instance_id while it is in state '$state'." >&2
      return 1
      ;;
  esac
}

awscli_run_ssm_script() {
  local instance_id="${1:?Instance ID is required}"
  local script_uri="${2:?S3 script URI is required}"
  local comment="${3:-Run Authify deployment script}"
  local remote_command parameters_file command_id invocation status attempt

  awscli_wait_for_ssm_instance "$instance_id"

  printf -v remote_command \
    'export AWS_REGION=%q AWS_DEFAULT_REGION=%q; cloud-init status --wait && aws s3 cp %q /tmp/authify-deploy.sh --region %q && bash /tmp/authify-deploy.sh' \
    "$REGION" "$REGION" "$script_uri" "$REGION"

  parameters_file="$(mktemp)"
  jq -n --arg command "$remote_command" '{commands: [$command]}' > "$parameters_file"
  if ! command_id="$(aws ssm send-command \
    --document-name AWS-RunShellScript \
    --instance-ids "$instance_id" \
    --comment "$comment" \
    --parameters "file://$parameters_file" \
    --query 'Command.CommandId' \
    --output text --region "$REGION")"; then
    rm -f "$parameters_file"
    return 1
  fi
  rm -f "$parameters_file"

  for ((attempt = 1; attempt <= 120; attempt++)); do
    invocation="$(aws ssm get-command-invocation \
      --command-id "$command_id" \
      --instance-id "$instance_id" \
      --output json --region "$REGION" 2>/dev/null || true)"
    status="$(jq -r '.Status // empty' <<< "$invocation" 2>/dev/null || true)"
    case "$status" in
      Success)
        echo "SSM deployment command succeeded on $instance_id."
        return 0
        ;;
      Failed|TimedOut|Cancelled|Cancelling)
        echo "SSM deployment command on $instance_id ended with status $status." >&2
        jq -r '.StandardErrorContent // empty' <<< "$invocation" >&2
        jq -r '.StandardOutputContent // empty' <<< "$invocation" >&2
        return 1
        ;;
    esac
    sleep 5
  done

  echo "SSM deployment command $command_id did not finish on $instance_id." >&2
  return 1
}
