### Sendsei Hem Charts Library

## Usage

The quickest way to use these charts is to simply copy an existing services infra folder.

For more information including developing a new chart type see: https://helm.sh/docs/topics/library_charts/

## Chart Types

-   k8sServiceChart - The most common Helm Chart used. This deploys a standard Kubernetes Service.
-   k8sCronJobChart - This chart can be used to deploy cronjobs to Kubernetes.
-   k8sDeploymentJobChart - This chart is used to run one off jobs at deployment time which terminate after being run once. Mainly used for migration or data update scripts.

## TEMPORARY WORKAROUND: Update all Chart Dependencies

There is a helper script which allows you to update all service Chart Dependencies with any changed made to a Chart Library.

This is a temproray workaround until we have a chart repository to host charts.
