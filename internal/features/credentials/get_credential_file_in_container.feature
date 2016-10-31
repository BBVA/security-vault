Feature:
  As a container
  In order to authenticate with external services
  I want to have my credentials somewhere

  @destroyContainers
  Scenario: Get a volume with my credentials
    Given a container "cred-test" configured         |
    When the container "cred-test" is started
    Then the container "cred-test" credentials will be the following
      | file         | content         |
      | /tmp/public | -BEGIN CERTIFICATE---- |
      | /tmp/private | -----BEGIN PRIVATE KEY-----|
      | /tmp/cacert     | -----BEGIN CERTIFICATE----- |