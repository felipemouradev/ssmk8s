SA=ssmk8s-sa
SECRET_NAME=$(kubectl get sa $SA -o jsonpath='{.secrets[0].name}')
TOKEN=$(kubectl get secret $SECRET_NAME -o jsonpath='{.data.token}' | base64 --decode)
CA_CRT=$(kubectl get secret $SECRET_NAME -o jsonpath='{.data.ca\.crt}')
APISERVER=$(kubectl config view --minify | grep server | cut -f 2- -d ":" | tr -d " ")

sed "s,CA_CRT,"$CA_CRT",g" ./kubeconfig_example > kubeconfig
sed -i -e "s,TOKEN,"$TOKEN",g" kubeconfig
sed -i -e "s,APISERVER,"$APISERVER",g" kubeconfig
sed -i -e "s,SA,"$SA",g" kubeconfig