# this is how to do this properly, but we use the null_resource for now
# https://medium.com/@marcincuber/amazon-eks-with-oidc-provider-iam-roles-for-kubernetes-services-accounts-59015d15cb0c
# resource "aws_iam_openid_connect_provider" "iam_oidc_provider" {
# url = aws_eks_cluster.cluster.identity.0.oidc.0.issuer
# client_id_list = ["sts.amazonaws.com"]
# thumbprint_list = [] # this is the nasty part
# }

resource "null_resource" "iam_oidc_provider" {
  provisioner "local-exec" {
    command = "eksctl utils associate-iam-oidc-provider --region ap-southeast-2 --cluster ${terraform.workspace} --approve"
  }
}

resource "aws_iam_policy" "alb_ingress_controller" {
  name   = "alb-ingress-controller-iam-policy-${terraform.workspace}"
  policy = file("${path.module}/policy.json")
}

resource "kubernetes_cluster_role" "alb_ingress_controller" {
  metadata {
    name = "alb-ingress-controller"

    labels = {
      "managed-by"             = "terraform"
      "cluster-name"           = terraform.workspace
      "app.kubernetes.io/name" = "alb-ingress-controller"
    }
  }

  rule {
    verbs      = ["create", "get", "list", "update", "watch", "patch"]
    api_groups = ["", "extensions"]
    resources  = ["configmaps", "endpoints", "events", "ingresses", "ingresses/status", "services"]
  }

  rule {
    verbs      = ["get", "list", "watch"]
    api_groups = ["", "extensions"]
    resources  = ["nodes", "pods", "secrets", "services", "namespaces"]
  }
}

resource "kubernetes_cluster_role_binding" "alb_ingress_controller" {
  metadata {
    name = "alb-ingress-controller"

    labels = {
      "managed-by"             = "terraform"
      "cluster-name"           = terraform.workspace
      "app.kubernetes.io/name" = "alb-ingress-controller"
    }
  }

  subject {
    kind      = "ServiceAccount"
    name      = "alb-ingress-controller"
    namespace = "kube-system"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "alb-ingress-controller"
  }
}

resource "null_resource" "iamserviceaccount" {
  provisioner "local-exec" {
    # this commmand automatically creates a kubernetes_service_account resource with it
    # this is AWS specific, but we're not multi cloud, so it's ok
    command = "eksctl create iamserviceaccount --region ap-southeast-2 --name alb-ingress-controller --namespace kube-system --cluster ${terraform.workspace} --attach-policy-arn arn:aws:iam::$(aws sts get-caller-identity | jq -r '.Account'):policy/alb-ingress-controller-iam-policy-${terraform.workspace} --override-existing-serviceaccounts --approve"
  }

  provisioner "local-exec" {
    when    = destroy
    command = "eksctl delete iamserviceaccount --region ap-southeast-2 --name alb-ingress-controller --namespace kube-system --cluster ${terraform.workspace}"
  }

  # alb_ingress_controller must be created before the iamserviceaccount
  depends_on = [aws_iam_policy.alb_ingress_controller]
}

resource "kubernetes_deployment" "alb_ingress_controller" {
  metadata {
    name = "alb-ingress-controller"
    namespace = "kube-system"
    labels = {
      "app.kubernetes.io/name" = "alb-ingress-controller"
    }
  }

  spec {
    selector {
      match_labels = {
        "app.kubernetes.io/name" = "alb-ingress-controller"
      }
    }

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name" = "alb-ingress-controller"
        }
      }

      spec {
        container {
          name  = "alb-ingress-controller"
          image = "docker.io/amazon/aws-alb-ingress-controller:v1.1.7"
          
          args = [
              "--ingress-class=alb",
              "--cluster-name=${terraform.workspace}"
          ]
        }

        service_account_name = "alb-ingress-controller"
        automount_service_account_token = "true"
      }
    }
  }
}
