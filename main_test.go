package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

var (
	credentials = []byte(`	
[sts-session]
aws_access_key_id = testAccessKeyID
aws_secret_access_key = testSecretKey
region = us-east-1
aws_session_token = testSessionToken

[iam-user]
aws_access_key_id = testAccessKeyID
aws_secret_access_key = testSecretKey
region = us-east-1

[derived]
role_arn = testRoleArn
source_profile = iam-user
role_session_name = derived

[bad-derived]
source_profile = iam-user

[bad-no-key]
aws_secret_access_key = testSecretKey

[bad-no-secret]
aws_access_key_id = testSecretKey
`)
)

func TestOktaSetup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS Profile Parser")
}

var _ = Describe("AWS Profile Parser", func() {
	Context("Happy Path", func() {
		It("should parser valid STS profile correctly", func() {
			awsProfile, err := ParseAWSProfile(credentials, "sts-session")
			Expect(err).To(Not(HaveOccurred()))

			Expect(awsProfile.AccessKeyID).To(Equal("testAccessKeyID"))
			Expect(awsProfile.SecretAccessKey).To(Equal("testSecretKey"))
			Expect(awsProfile.SessionToken).To(Equal("testSessionToken"))
			Expect(awsProfile.Region).To(Equal("us-east-1"))

			Expect(awsProfile.RoleArn).To(BeEmpty())
			Expect(awsProfile.RoleSessionName).To(BeEmpty())
			Expect(awsProfile.SourceProfile).To(BeEmpty())
		})

		It("should parser valid IAM user profile correctly", func() {
			awsProfile, err := ParseAWSProfile(credentials, "iam-user")
			Expect(err).To(Not(HaveOccurred()))
			Expect(awsProfile.AccessKeyID).To(Equal("testAccessKeyID"))
			Expect(awsProfile.SecretAccessKey).To(Equal("testSecretKey"))
			Expect(awsProfile.SessionToken).To(BeEmpty())
			Expect(awsProfile.Region).To(Equal("us-east-1"))

			Expect(awsProfile.RoleArn).To(BeEmpty())
			Expect(awsProfile.RoleSessionName).To(BeEmpty())
			Expect(awsProfile.SourceProfile).To(BeEmpty())
		})

		It("should parser valid derived profile correctly", func() {
			awsProfile, err := ParseAWSProfile(credentials, "derived")
			Expect(err).To(Not(HaveOccurred()))

			Expect(awsProfile.RoleSessionName).To(Equal("derived"))
			Expect(awsProfile.RoleArn).To(Equal("testRoleArn"))
			Expect(awsProfile.SourceProfile).To(Equal("iam-user"))

			Expect(awsProfile.AccessKeyID).To(BeEmpty())
			Expect(awsProfile.SecretAccessKey).To(BeEmpty())
			Expect(awsProfile.SessionToken).To(BeEmpty())
			Expect(awsProfile.Region).To(BeEmpty())
		})
	})

	Context("Error Path", func() {
		It("should return error for invalid profile name", func() {
			_, err := ParseAWSProfile(credentials, "non-existent")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no profile found with name: `non-existent`"))
		})

		It("should return error for invalid derived profile", func() {
			_, err := ParseAWSProfile(credentials, "bad-derived")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no `role_arn` found in profile"))
		})

		It("should return error for invalid profile with no access key", func() {
			_, err := ParseAWSProfile(credentials, "bad-no-key")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no `aws_access_key_id` found in profile"))
		})

		It("should return error for invalid profile with no secret key", func() {
			_, err := ParseAWSProfile(credentials, "bad-no-secret")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no `aws_secret_access_key` found in profile"))
		})
	})
})
