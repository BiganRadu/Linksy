load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_dependencies():
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2",
        importpath = "github.com/aws/aws-sdk-go-v2",
        sum = "h1:mJoei2CxPutQVxaATCzDUjcZEjVRdpsiiXi2o38yqWM=",
        version = "v1.36.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_aws_protocol_eventstream",
        importpath = "github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream",
        sum = "h1:zAybnyUQXIZ5mok5Jqwlf58/TFE7uvd3IAsa1aF9cXs=",
        version = "v1.6.10",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_config",
        importpath = "github.com/aws/aws-sdk-go-v2/config",
        sum = "h1:yNjgjiGBp4GgaJrGythyBXg2wAs+Im9fSWIUwvi1CAc=",
        version = "v1.29.10",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_credentials",
        importpath = "github.com/aws/aws-sdk-go-v2/credentials",
        sum = "h1:rv1V3kIJ14pdmTu01hwcMJ0WAERensSiD9rEWEBb1Tk=",
        version = "v1.17.63",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_feature_ec2_imds",
        importpath = "github.com/aws/aws-sdk-go-v2/feature/ec2/imds",
        sum = "h1:x793wxmUWVDhshP8WW2mlnXuFrO4cOd3HLBroh1paFw=",
        version = "v1.16.30",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_internal_configsources",
        importpath = "github.com/aws/aws-sdk-go-v2/internal/configsources",
        sum = "h1:ZK5jHhnrioRkUNOc+hOgQKlUL5JeC3S6JgLxtQ+Rm0Q=",
        version = "v1.3.34",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_internal_endpoints_v2",
        importpath = "github.com/aws/aws-sdk-go-v2/internal/endpoints/v2",
        sum = "h1:SZwFm17ZUNNg5Np0ioo/gq8Mn6u9w19Mri8DnJ15Jf0=",
        version = "v2.6.34",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_internal_ini",
        importpath = "github.com/aws/aws-sdk-go-v2/internal/ini",
        sum = "h1:bIqFDwgGXXN1Kpp99pDOdKMTTb5d2KyU5X/BZxjOkRo=",
        version = "v1.8.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_internal_v4a",
        importpath = "github.com/aws/aws-sdk-go-v2/internal/v4a",
        sum = "h1:ZNTqv4nIdE/DiBfUUfXcLZ/Spcuz+RjeziUtNJackkM=",
        version = "v1.3.34",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_internal_accept_encoding",
        importpath = "github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding",
        sum = "h1:eAh2A4b5IzM/lum78bZ590jy36+d/aFLgKF/4Vd1xPE=",
        version = "v1.12.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_internal_checksum",
        importpath = "github.com/aws/aws-sdk-go-v2/service/internal/checksum",
        sum = "h1:lguz0bmOoGzozP9XfRJR1QIayEYo+2vP/No3OfLF0pU=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_internal_presigned_url",
        importpath = "github.com/aws/aws-sdk-go-v2/service/internal/presigned-url",
        sum = "h1:dM9/92u2F1JbDaGooxTq18wmmFzbJRfXfVfy96/1CXM=",
        version = "v1.12.15",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_internal_s3shared",
        importpath = "github.com/aws/aws-sdk-go-v2/service/internal/s3shared",
        sum = "h1:moLQUoVq91LiqT1nbvzDukyqAlCv89ZmwaHw/ZFlFZg=",
        version = "v1.18.15",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_s3",
        importpath = "github.com/aws/aws-sdk-go-v2/service/s3",
        sum = "h1:jIiopHEV22b4yQP2q36Y0OmwLbsxNWdWwfZRR5QRRO4=",
        version = "v1.78.2",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_sso",
        importpath = "github.com/aws/aws-sdk-go-v2/service/sso",
        sum = "h1:8JdC7Gr9NROg1Rusk25IcZeTO59zLxsKgE0gkh5O6h0=",
        version = "v1.25.1",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_ssooidc",
        importpath = "github.com/aws/aws-sdk-go-v2/service/ssooidc",
        sum = "h1:wK8O+j2dOolmpNVY1EWIbLgxrGCHJKVPm08Hv/u80M8=",
        version = "v1.29.2",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_sts",
        importpath = "github.com/aws/aws-sdk-go-v2/service/sts",
        sum = "h1:PZV5W8yk4OtH1JAuhV2PXwwO9v5G5Aoj+eMCn4T+1Kc=",
        version = "v1.33.17",
    )
    go_repository(
        name = "com_github_aws_smithy_go",
        importpath = "github.com/aws/smithy-go",
        sum = "h1:6D9hW43xKFrRx/tXXfAlIZc4JI+yQe6snnWcQyxSyLQ=",
        version = "v1.22.2",
    )
    go_repository(
        name = "com_github_skip2_go_qrcode",
        importpath = "github.com/skip2/go-qrcode",
        sum = "h1:MRM5ITcdelLK2j1vwZ3Je0FKVCfqOLp5zO6trqMLYs0=",
        version = "v0.0.0-20200617195104-da1b6568686e",
    )
