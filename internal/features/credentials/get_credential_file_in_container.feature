Feature:
  As a container
  In order to authenticate with external services
  I want to have my credentials somewhere

  @destroyContainers
  Scenario: Get a volume with my credentials
    Given a container "cred-test" configured with the following volume driver options:
      | volume_driver | host_mount_point | container_mount_point |
      | Vault         | test/mountpoint  | /secret               |
    When the container "cred-test" is started
    Then the container "cred-test" credentials will be the following
      | file               | content          |
      | /secret/cert | certificadooorr |