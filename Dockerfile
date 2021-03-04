FROM golang

WORKDIR /go/src/github.com/billtrust/terraform-provider-looker

RUN apt-get update && apt-get install unzip

RUN wget https://releases.hashicorp.com/terraform/0.14.5/terraform_0.14.5_linux_amd64.zip && \
    unzip terraform_0.14.5_linux_amd64.zip && \
    chmod +x terraform && \
    mv terraform /usr/local/bin

#COPY ./ .
RUN go get github.com/gruntwork-io/terratest/modules/terraform
