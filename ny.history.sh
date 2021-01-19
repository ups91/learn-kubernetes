helm dependency list ./helm-chart/; helm dependency update ./helm-chart/; helm dependency build ./helm-chart/;

helm install my-echo ./helm-chart/ --namespace echo-namespace --set hosts.host=test.com --create-namespace