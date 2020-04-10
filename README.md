# AWS Profile Parser

AWS profile parser reads an [AWS Credentials File](https://docs.aws.amazon.com/sdk-for-php/v3/developer-guide/guide_credentials_profiles.html) and echos it in JSON format. 

### Usage 

**Good**:

```json 
./aws-profile-parser -c ../test/.aws/credentials -p test
{
  "AWS_ACCESS_KEY_ID": "ASIA2974r973492",
  "AWS_SECRET_ACCESS_KEY": "FVHxnIS/Y24308oufofuf0",
  "AWS_SESSION_TOKEN": "AOUEdlsjfewur0wrujo/Hv1Q==",
  "AWS_REGION": "us-east-1"
}
``` 

```json 
./aws-profile-parser -c ../test/.aws/credentials -p test2
{
  "AWS_ROLE_ARN": "arn:aws:iam::00000000000:role/ReadOnly",
  "AWS_SOURCE_PROFILE": "test"
}
```

**Bad**:

```bash
 ./aws-profile-parser -c ../endgame-sre-infrastructure/tools/okta-login/.aws/credentials -p test-bad
Error: invalid AWS profile [test-bad]: no `source_profile` found in profile
Usage:
  aw-profile-parser [flags]

Flags:
  -c, --credentials string   AWS credentials file
  -h, --help                 help for aw-profile-parser
  -p, --profile string       AWS profile name (default "default")
  -v, --version              version for aw-profile-parser

FATA[0000] Error executing command: invalid AWS profile [test-bad]: no `source_profile` found in profile
```

## Build & Release

To build this project simply run:

```bash
go build .
```

This project uses [GoReleaser](https://goreleaser.com) to publish releases. This is currently not integrated into CI and done manually. 