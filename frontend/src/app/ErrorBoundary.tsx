import { Component, type PropsWithChildren } from 'react';
import { PageState } from '../shared/ui/PageState';

export class ErrorBoundary extends Component<PropsWithChildren, { hasError: boolean }> {
  state = { hasError: false };

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  render() {
    if (this.state.hasError) {
      return (
        <PageState
          title="Something went wrong"
          description="Please reload the page to try again."
          onRetry={() => window.location.reload()}
        />
      );
    }
    return this.props.children;
  }
}
