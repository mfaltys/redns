 #!/bin/bash
 echo $1 | gpg \
	--passphrase-fd 0 \
	--batch --yes \
	--no-default-keyring --armor \
	--secret-keyring ./unixvoid.sec --keyring ./unixvoid.pub \
	--output doic-latest-linux-amd64.aci.asc \
	--detach-sig doic-latest-linux-amd64.aci
