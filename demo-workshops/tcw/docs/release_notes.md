# Release Notes

## What's Changed

### v2.12

1. (fbc8dff) Fixes #18 where automatic updates prevent ansible roles from running 

### v2.11
1. (4d8014c) Fixes race condition issue that manifests in how_many > 1 environments

### v2.10
1. This repo now supports building on Windows systems

### v2.09
1. (9ab28a0) removed environment echo
2. resolved delphix/dx-workshops#7

### v2.08

1. (6fd064d) Collapsed build steps for LabAlchemy-compliant systems
2. (b4963fe) Resolved https://github.com/delphix/dx-workshops/issues/3

### v2.07

1. (5f6f8da) workaround for GUAC-414 as found https://github.com/cisagov/pca-gophish-composition-packer/pull/1

### v2.06

1. (f75e1d3) Workshop environments now support multiple students (HOW_MANY > 1).

### v2.05

1. (50e2c0c) credentials added to desktop wallpaper
2. (9290da7) show bookmark bar
3. (ee57992) added masking engine bookmark

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

1. None, at this time.