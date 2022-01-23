FROM    golang:1.17.6-alpine

# Setup
EXPOSE  9000
ENV     CONFIG prod # Set local or prod

# Copy project file and Go modules
WORKDIR /app
COPY    . ./
RUN     go mod download
RUN     go mod verify

# Build
RUN     go build -o /wms_slave
CMD     [ "/wms_slave" ]