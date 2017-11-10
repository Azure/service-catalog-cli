# Service Catalog Demonstration Walkthrough

This document assumes that you've installed Service Catalog onto your cluster.
If you haven't, please see [install.md](./install.md).

All commands in this document assume that you're operating out of the root
of this repository.

# Step 1 - Installing the UPS Broker Server

Since the Service Catalog provides a Kubernetes-native interface to an 
[Open Service Broker API](https://www.openservicebrokerapi.org/) compatible broker
server, we'll need to install one in order to proceed with a demo.

In this repository, there's a simple, "dummy" server called the User Provided 
Service (UPS) broker. The codebase for that broker is
[here](https://github.com/kubernetes-incubator/service-catalog/tree/master/contrib/pkg/broker/user_provided/controller).

We're going to deploy the UPS broker to our Kubernetes cluster before 
proceeding, and we'll do so with the UPS helm chart. You can find details about 
that chart in the chart's 
[README](https://github.com/kubernetes-incubator/service-catalog/blob/master/charts/ups-broker/README.md).

Otherwise, to install with sensible defaults, run the following command:

```console
helm install charts/ups-broker --name ups-broker --namespace ups-broker
```

# Step 2 - Creating a `ClusterServiceBroker` Resource

Because we haven't created any resources in the service-catalog API server yet,
`kubectl get` will return an empty list of resources.

```console
kubectl get clusterservicebrokers,clusterserviceclasses,serviceinstances,servicebindings
No resources found.
```

We'll register a broker server with the catalog by creating a new
[`ClusterServiceBroker`](../contrib/examples/walkthrough/ups-broker.yaml) resource.
Do so with the following command:

```console
kubectl create -f contrib/examples/walkthrough/ups-broker.yaml
```

The output of that command should be the following:

```console
servicebroker "ups-broker" created
```

When we create this `ClusterServiceBroker` resource, the service catalog controller responds
by querying the broker server to see what services it offers and creates a
`ClusterServiceClass` for each.

We can check the status of the broker using `kubectl get`:

```console
kubectl get clusterservicebrokers ups-broker -o yaml
```

We should see something like:

```yaml
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ClusterServiceBroker
metadata:
  creationTimestamp: 2017-03-03T04:11:17Z
  finalizers:
  - kubernetes-incubator/service-catalog
  name: ups-broker
  resourceVersion: "6"
  selfLink: /apis/servicecatalog.k8s.io/v1beta1/clusterservicebrokers/ups-broker
  uid: 72fa629b-ffc7-11e6-b111-0242ac110005
spec:
  url: http://ups-broker-ups-broker.ups-broker.svc.cluster.local
status:
  conditions:
  - message: Successfully fetched catalog entries from broker.
    reason: FetchedCatalog
    status: "True"
    type: Ready
```

Notice that the `status` field has been set to reflect that the broker server's
catalog of service offerings has been successfully added to our cluster's
service catalog.

# Step 3 - Viewing `ClusterServiceClass`es

The controller created a `ClusterServiceClass` for each service that the UPS broker
provides. We can view the `ClusterServiceClass` resources available in the cluster by
executing:

```console
$ kubectl get clusterserviceclasses -o=custom-columns=NAME:.metadata.name,EXTERNAL\ NAME:.spec.externalName
```

We should see something like:

```console
NAME                                   EXTERNAL NAME
4f6e6cf6-ffdd-425f-a2c7-3c9258ad2468   user-provided-service
```

**NOTE:** The above command uses a custom set of columns.  The `NAME` field is
the Kubernetes name of the ClusterServiceClass and the `EXTERNAL NAME` field is the
human-readable name for the service that the broker returns.

The UPS broker provides a service with the external name
`user-provided-service`. Run the following command to see the details of this
offering:

```console
kubectl get clusterserviceclasses 4f6e6cf6-ffdd-425f-a2c7-3c9258ad2468 -o yaml
```

We should see something like:

```yaml
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ClusterServiceClass
metadata:
  creationTimestamp: 2017-03-03T04:11:17Z
  name: user-provided-service
  resourceVersion: "7"
  selfLink: /apis/servicecatalog.k8s.io/v1beta1/clusterserviceclasses/user-provided-service
  uid: 72fef5ce-ffc7-11e6-b111-0242ac110005
brokerName: ups-broker
externalID: 4F6E6CF6-FFDD-425F-A2C7-3C9258AD2468
bindable: false
planUpdatable: false
plans:
- name: default
  free: true
  externalID: 86064792-7ea2-467b-af93-ac9694d96d52
```

# Step 4 - Creating a New `ServiceInstance`

Now that a `ClusterServiceClass` named `user-provided-service` exists within our
cluster's service catalog, we can create a `ServiceInstance` that points to
it.

Unlike `ClusterServiceBroker` and `ClusterServiceClass` resources, `ServiceInstance` 
resources must be namespaced, so we'll need to create a namespace to start.
Do so with this command:

```console
kubectl create namespace test-ns
```

Then, create the `ServiceInstance`:

```console
kubectl create -f contrib/examples/walkthrough/ups-instance.yaml
```

That operation should output:

```console
serviceinstance "ups-instance" created
```

After the `ServiceInstance` is created, the service catalog controller will 
communicate with the appropriate broker server to initiate provisioning. 
Check the status of that process with this command:

```console
kubectl get serviceinstances -n test-ns ups-instance -o yaml
```

We should see something like:

```yaml
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ServiceInstance
metadata:
  creationTimestamp: 2017-10-02T14:50:28Z
  finalizers:
  - kubernetes-incubator/service-catalog
  generation: 1
  name: ups-instance
  namespace: test-ns
  resourceVersion: "12"
  selfLink: /apis/servicecatalog.k8s.io/v1beta1/namespaces/test-ns/serviceinstances/ups-instance
  uid: 07ecf19d-a781-11e7-8b18-0242ac110005
spec:
  externalID: 7f2c176a-ae67-4b5e-a826-58591d85a1d7
  clusterServiceClassExternalName: user-provided-service
  clusterServicePlanExternalName: default
  parameters:
    credentials:
      param-1: value-1
      param-2: value-2
status:
  asyncOpInProgress: false
  conditions:
  - lastTransitionTime: 2017-10-02T14:50:28Z
    message: The instance was provisioned successfully
    reason: ProvisionedSuccessfully
    status: "True"
    type: Ready
  externalProperties:
    clusterServicePlanExternalName: default
    parameterChecksum: e65c764db8429f9afef45f1e8f71bcbf9fdbe9a13306b86fd5dcc3c5d11e5dd3
    parameters:
      credentials:
        param-1: value-1
        param-2: value-2
  orphanMitigationInProgress: false
  reconciledGeneration: 1

```

# Step 5 - Requesting a `ServiceBinding` to use the `ServiceInstance`

Now that our `ServiceInstance` has been created, we can bind to it. To accomplish this,
we'll create a `ServiceBinding` resource. Do so with the following
command:

```console
kubectl create -f contrib/examples/walkthrough/ups-binding.yaml
```

That command should output:

```console
servicebinding "ups-binding" created
```

After the `ServiceBinding` resource is created, the service catalog controller will
communicate with the appropriate broker server to initiate binding. Generally,
this will cause the broker server to create and issue credentials that the
service catalog controller will insert into a Kubernetes `Secret`. We can check
the status of this process like so:

```console
kubectl get servicebindings -n test-ns ups-binding -o yaml
```

We should see something like:

```yaml
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ServiceBinding
metadata:
  creationTimestamp: 2017-03-07T01:44:36Z
  finalizers:
  - kubernetes-incubator/service-catalog
  name: ups-binding
  namespace: test-ns
  resourceVersion: "29"
  selfLink: /apis/servicecatalog.k8s.io/v1beta1/namespaces/test-ns/servicebindings/ups-binding
  uid: 9eb2cdce-02d7-11e7-8edb-0242ac110005
spec:
  instanceRef:
    name: ups-instance
  externalID: b041db94-a5a0-41a2-87ae-1025ba760918
  secretName: ups-binding
status:
  conditions:
  - lastTransitionTime: 2017-03-03T01:44:37Z
    message: Injected bind result
    reason: InjectedBindResult
    status: "True"
    type: Ready
```

Notice that the status has a `Ready` condition set.  This means our binding is
ready to use!  If we look at the `Secret`s in our `test-ns` namespace, we should
see a new one:

```console
kubectl get secrets -n test-ns
NAME                              TYPE                                  DATA      AGE
default-token-3k61z               kubernetes.io/service-account-token   3         29m
ups-binding                       Opaque                                2         1m
```

Notice that a new `Secret` named `ups-binding` has been created.

# Step 6 - Deleting the `ServiceBinding`

Now, let's unbind from the instance. To do this, we simply *delete* the
`ServiceBinding` resource that we previously created:

```console
kubectl delete -n test-ns servicebindings ups-binding
```

After the deletion is complete, we should see that the `Secret` is gone:

```console
kubectl get secrets -n test-ns
NAME                  TYPE                                  DATA      AGE
default-token-3k61z   kubernetes.io/service-account-token   3         30m
```

# Step 7 - Deleting the `ServiceInstance`

Now, we can deprovision the instance. To do this, we simply *delete* the
`ServiceInstance` resource that we previously created:

```console
kubectl delete -n test-ns serviceinstances ups-instance
```

# Step 8 - Deleting the `ClusterServiceBroker`

Next, we should remove the `ClusterServiceBroker` resource. This tells the service
catalog to remove the broker's services from the catalog. Do so with this
command:

```console
kubectl delete clusterservicebrokers ups-broker
```

We should then see that all the `ClusterServiceClass` resources that came from that
broker have also been deleted:

```console
kubectl get clusterserviceclasses
No resources found.
```

# Step 9 - Final Cleanup

## Cleaning up the UPS Service Broker Server

To clean up, delete the helm deployment:

```console
helm delete --purge ups-broker
```

Then, delete all the namespaces we created:

```console
kubectl delete ns test-ns ups-broker
```
## Cleaning up the Service Catalog

Delete the helm deployment and the namespace:

```console
helm delete --purge catalog
kubectl delete ns catalog
```

# Troubleshooting

## Firewall rules

If you are using Google Cloud Platform, you may need to run the following
commands to setup proper firewall rules to allow your traffic get in.

```console
gcloud compute firewall-rules create allow-service-catalog-secure --allow tcp:30443 --description "Allow incoming traffic on 30443 port."
```
