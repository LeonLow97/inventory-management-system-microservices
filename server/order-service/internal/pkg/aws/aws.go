package aws

// const DEFAULT_REGION = "ap-southeast-1"

// // NewSession creates a new session with aws
// func NewSession(config config.Config) (*session.Session, error) {
// 	if config.AWS.Region == "" {
// 		config.AWS.Region = DEFAULT_REGION
// 	}

// 	awsConfig := aws.Config{
// 		Credentials:      credentials.NewStaticCredentials(config.AWS.Credentials.ID, config.AWS.Credentials.Secret, config.AWS.Credentials.Token),
// 		Region:           aws.String(config.AWS.Region),
// 		Endpoint:         aws.String(config.AWS.Endpoint),
// 		S3ForcePathStyle: aws.Bool(config.AWS.S3ForcePathStyle),
// 	}

// 	return session.NewSessionWithOptions(
// 		session.Options{
// 			Config: awsConfig,
// 		},
// 	)
// }
