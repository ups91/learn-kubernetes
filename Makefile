NAME_SPACE := $(shell awk '$$1=="name:"{print $$2;exit(0)}' ./k8s_conf/namespace.yaml)
CLUSTER_NAME := echo_cluster

build:
	cd ./echo   ; go build  -o ./echo.so
	cd ./timer  ; go build  -o ./timer.so
	cd ./colorer; go build  -o ./colorer.so

.PHONY:docker
docker:
	rm -rf ./echo/linux
	cd ./echo;\
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./linux/echo_linux.so ;\
	mkdir ./linux/d ;\
	chmod 777 ./linux/d ;\
	docker build -f ./Dockerfile -t ups91/kube-test-pods:echo --no-cache --force-rm ./linux/ ;\
	docker push  ups91/kube-test-pods:echo
	
	rm -rf ./timer/linux
	cd ./timer;\
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./linux/timer_linux.so ;\
	mkdir ./linux/d ;\
	chmod 777 ./linux/d ;\
	docker build -f ./Dockerfile -t ups91/kube-test-pods:timer --no-cache --force-rm ./linux/ ;\
	docker push  ups91/kube-test-pods:timer
	
	rm -rf ./colorer/linux
	cd ./colorer;\
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./linux/colorer_linux.so ;\
	mkdir ./linux/d ;\
	chmod 777 ./linux/d ;\
	docker build -f ./Dockerfile -t ups91/kube-test-pods:colorer --no-cache --force-rm ./linux/ ;\
	docker push  ups91/kube-test-pods:colorer

docker-login:
	#docker login

.PHONY:clean
clean:
	rm -rf ./linux&
	rm -rf ./main&
	rm -rf ./echo/linux
	rm -rf ./timer/linux
	rm -rf ./colorer/linux
	
	go clean -cache -modcache -i -r&

	docker image prune -a -f --filter "label=echo"
	docker image prune -a -f --filter "label=timer"
	docker image prune -a -f --filter "label=colorer"
	#docker stop echo_contemner
	#docker rm -f echo_container

.PHONY:k8s-up
k8s-up:
	kubectl apply -f ./k8s_conf/namespace.yaml
	kubectl apply -f ./k8s_conf/azure-PVC.yaml
	#
	#find out Resource group
	#az aks show --resource-group echo_cluster --name echo_cluster --query nodeResourceGroup -o tsv
	#create static IPs
	#az network public-ip create --resource-group MC_echo_cluster_echo_cluster_westeurope --name echo_PublicIP --sku Standard --allocation-method static --query publicIp.ipAddress -o tsv



	
	kubectl apply -f ./k8s_conf/deploy-Set-Timer.yaml --namespace=$(NAME_SPACE)
	kubectl apply -f ./k8s_conf/service-timer.yaml    --namespace=$(NAME_SPACE)

	kubectl apply -f ./k8s_conf/deploy-Set-Colorer.yaml --namespace=$(NAME_SPACE)
	kubectl apply -f ./k8s_conf/service-colorer.yaml    --namespace=$(NAME_SPACE)

	kubectl apply -f ./k8s_conf/deploy-Set-Echo.yaml  --namespace=$(NAME_SPACE)
	kubectl apply -f ./k8s_conf/service-echo.yaml     --namespace=$(NAME_SPACE)



	##################VAULT ##############
	#az keyvault create --resource-group echo_cluster2 --name echo-vault
	#az keyvault secret set \
	#	--vault-name echo-vault2 \
	#	--name SecretPassSecretKeyssword \
	#	--value SuperHiddenString
	##Debug
	#az keyvault secret show --name "SecretPassSecretKeyssword" --vault-name "echo-vault2"
	#####################################

.PHONY:k8s-rm
k8s-rm:
	kubectl delete -f ./k8s_conf/namespace.yaml --namespace=$(NAME_SPACE)
	#az keyvault delete --resource-group echo_cluster --name echo-vault
	#az keyvault purge --name echo-vault


.PHONY:requirements
requirements:
	# Add the ingress-nginx repository
	# Use Helm to deploy an NGINX ingress controller
	kubectl apply -f ./k8s_conf/namespace-requirements.yaml
	#
	## helm repo add csi-secrets-store-provider-azure https://raw.githubusercontent.com/Azure/secrets-store-csi-driver-provider-azure/master/charts
	## helm install csi-secrets-store-provider-azure/csi-secrets-store-provider-azure --generate-name

	kubectl apply -f ./k8s_conf/secret.yaml           --namespace=echo-requirements
	

	#helm install nginx-ing ingress-nginx/ingress-nginx --set rbac.create=true --namespace=echo-requirements
	kubectl apply -f ./k8s_conf/ingress-ngx.yaml      --namespace=echo-requirements
	#kubectl get validatingwebhookconfigurations  --all-namespaces
	#kubectl delete validatingwebhookconfigurations nginx-ing-ingress-nginx-admission-NAME

	helm repo add csi-secrets-store-provider-azure https://raw.githubusercontent.com/Azure/secrets-store-csi-driver-provider-azure/master/charts
	helm repo add secrets-store-csi-driver https://raw.githubusercontent.com/kubernetes-sigs/secrets-store-csi-driver/master/charts

	#helm install csi-secrets-store-provider-azure/csi-secrets-store-provider-azure --generate-name
	#helm install csi-secrets-store secrets-store-csi-driver/secrets-store-csi-driver --generate-name
	#
	kubectl apply -f https://raw.githubusercontent.com/Azure/secrets-store-csi-driver-provider-azure/master/deployment/provider-azure-installer.yaml --namespace=$(NAME_SPACE)


	

