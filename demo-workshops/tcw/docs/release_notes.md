# Release Notes

## What's Changed

### v2.04

1. Renamed Github Repo
2. Delphix Engine >=5.3.3 required
3. Datical 2019.2.2.6029 required
4. Upgraded to Terraform 12.5
5. Reorganized repo structure
6. EC2 instances inherit AMI tags
7. Significant improvements to the cleanup function
8. Git repos now owned by git user
9.  Jenkinsfile library now leverages docker container for all functions
10. DELPHIX_ADMIN_PASSWORD Password must be between 6-12 characters and contain 1 digit, 1 uppercase alphabet character, and 1 special character
11. Source for all Golang binaries used in workshops included in /go
12. All Golang binaries compiled at build time.
13. Engine setup scripts retry on 502 ( commonly occurs on initial boot )
14. Changed AMI names to be easier to work with
15. Eliminated extra sync+refresh step on workshop start (post deploy)
16. Added "Shutdown VDB's" shortcut to jumpbox desktop to facilitate clean shutdown of workshop
17. Parallelized some startup processes for fast boot time
18. Upgraded Guacamole to v1.0.0
19. Added Delphix customizations to Guacamole configuration
20. Added dev_mode flag to deploy - This enables public IP's for all workshop EC2 instances and configures the firewall for inbound access.
21. executing `docker-compose run tcw` without arguments prints the help menu.

## Known Issues

1. Frequent disconnect errors from Guacamole. Upstream Dependency [GUACAMOLE-414](https://issues.apache.org/jira/browse/GUACAMOLE-414)