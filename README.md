# IP Country Allowlist Service

## Overview

This project is a Go-based microservice that checks whether a given IP address originates from an **allowed country**. It uses MaxMind’s **GeoLite2 Country** database to map IP addresses to country codes, and a MySQL database to store a whitelist (allowlist) of allowed country ISO codes. The service exposes an HTTP API (and a gRPC API) for validating IP addresses. If the IP’s country is in the allowlist, the service indicates it’s allowed; otherwise it responds that the IP is not allowed.

## Features

- **IP Geolocation:** Uses the MaxMind GeoLite2 Country database to efficiently determine the country for any IP address (IPv4 or IPv6).
- **Country Allowlist:** Enforces access based on an allowlist of country codes stored in MySQL. All requests share the same global allowlist (no per-user or per-client differentiation at this time).
- **HTTP and gRPC APIs:** Provides a RESTful HTTP endpoint (`POST /check`) and a gRPC service (running on a separate port) to check IP addresses. This allows integration in different environments or with other microservices.
- **Fast Performance:** Loads the GeoLite2 data locally at startup for quick country lookups. Only a lightweight database query is performed on each request to verify the country code, keeping response times low.
- **Graceful Shutdown:** Handles OS signals for shutdown, allowing the service to terminate cleanly (closing server sockets and database connections).

## Architecture

This deployment consists of two main components: the **IP check service** and the **MySQL database**. The IP check service is a stateless microservice (containerized Go application) and MySQL acts as the stateful backend for configuration data (the allowed countries list).

When a request hits the IP check service, the following happens:

1. Extract IP Address
   When a request is received, the service extracts the IP address from it.

2. Fetch Allowed Country Codes
   The service queries a MySQL database to retrieve the list of allowed country codes from the whitelisted_tb table. These are stored as ISO country codes.

3. gRPC Verification
   The service then sends the extracted IP address and the list of allowed country codes to a gRPC API.
   The gRPC service performs the country lookup for the IP and determines whether the IP belongs to an allowed country.

4. Response Handling

If the gRPC service returns allowed, the service responds with HTTP 200 (OK).

If not allowed, the service responds with HTTP 417 (Expectation Failed).

Both the IP check service and MySQL database are deployed as separate containers (or Pods in Kubernetes) and communicate over the network. In the provided Kubernetes setup, they are in the same namespace and the IP check service uses the DNS name of the MySQL service to connect. The MySQL database persists the allowlist in a volume, so that allowed countries remain saved across restarts.

![Architecture Diagram: IP Check Service and MySQL]

- Diagram:
  The IP Check Service receives an IP address via HTTP, queries MySQL to retrieve the list of allowed country codes, and then sends the IP address and the allowlist to a gRPC service. The gRPC service determines the IP’s country using the GeoLite2 database and checks if it’s in the allowed list. If allowed, the service returns HTTP 200 OK; otherwise, it returns HTTP 417 Expectation Failed.

## Setup and Deployment

You can run this service either with Docker (for local testing) or deploy it to a Kubernetes cluster and docker-composer . Below are instructions for both methods.

### Running with Docker (Standalone)

# ```bash

git clone https://github.com/npx-json/assessment.git
cd assessment
docker compose up -d
or
minikube start
kubectl apply -f ipcheck.yaml

### API Usage

Once the service is running, you can use the HTTP API to check IP addresses. The service listens for HTTP GET requests at the /check endpoint. You should provide the IP address to check as a query parameter. For example:
Allow check (example): Suppose you have added “US” to the allowed countries. Checking an IP from the United States should return a success:
curl "http://localhost:8080/check"

If "US" is in the allowlist, the service will respond with HTTP 200 OK (allowed). The response body may be empty (the significance is in the status code).
Deny check (example): If you check an IP that is not from an allowed country (or if the allowlist is empty), you will get a failure response.

## GeoLite2 Database Updates

The GeoLite2 data file (GeoLite2-Country.mmdb) is included in the image and used for IP-to-country lookup. Automatic updates for this database are currently disabled. The code has functionality to update the GeoLite2 database on a schedule (it calls server.StartAutoUpdate() in the geoip module), but this feature was commented out because it requires a MaxMind license key to download updates. Currently, the service will continue using the static database file that it shipped with. Over time, this database will become outdated (new IP ranges might not be recognized or might map to old assignments). To update the GeoLite2 database manually: obtain a new GeoLite2 Country database (in .mmdb format) from MaxMind. You can sign up for a free MaxMind account and get a license key, then download the latest GeoLite2-Country.mmdb. Replace the file used by the service:
In Docker, you could build a new image with the updated .mmdb file, or mount the new file into the container and set GEOLITE_DB_PATH to the mount location.
In Kubernetes, you could create a ConfigMap or volume with the new DB file and mount it over the old one, or build a new image version. After replacing the file, restart the pod to load the new data.
If you wanted to enable automatic updates in code, you would need to provide the MaxMind license key to the application (likely via an environment variable or configuration) and uncomment/implement the update logic. As of now, keeping the database updated is a manual maintenance task.

## Limitations and Future Improvements

Global Allowlist Only: The current implementation uses a single, global allowlist of country codes. There is no support for different allowlists per customer or user group. In a multi-tenant scenario, every request is treated the same. A future improvement could introduce a way to have multiple allowlists or to scope the check by an additional parameter (e.g., a customer ID) and have a corresponding database schema to match IP country against that customer’s allowed countries.
Manual Allowlist Management: There is no API to add or remove countries from the allowed list at runtime. Management of the allowlist is manual – you must insert or delete rows in the MySQL table directly. In a production environment, you might build an admin interface or endpoint to manage this, or integrate with a config management system. For now, updating the allowlist likely involves running SQL commands or migrations (as shown in the setup steps).

## Dependency on Database Availability:

The microservice depends on MySQL being available each time it processes a request. On startup, if the database is unreachable, the service will fail to start (and in Kubernetes, it will crash loop until the DB comes up). There is no retry logic within the application for database queries – it assumes the DB is up. This is usually fine in a Kubernetes deployment (the liveness/readiness probes and container orchestration handle ordering), but it’s worth noting that a database outage will immediately affect the service’s ability to allow any requests. Using a lightweight cache for the allowlist (with periodic refresh) could be an enhancement to tolerate brief DB outages or reduce latency.
HTTP 417 Usage: The service uses HTTP status code 417 (“Expectation Failed”) to indicate a disallowed IP. This is a somewhat uncommon choice for an authorization failure (HTTP 403 Forbidden might be more typical). Clients calling the API need to be aware that 417 specifically means “blocked by allowlist policy” in this context. This is a design decision inherited from the original implementation. It works, but developers integrating the service should handle this status accordingly.

## GeoLite2 Update Process:

As described, the automatic update feature for the GeoLite2 database is not active. This means the accuracy of IP-to-country mapping will degrade over time until the database is updated manually. In a future iteration, providing a secure way to supply the MaxMind license key (e.g., via Kubernetes Secret or environment variable) and enabling periodic updates would make maintenance easier.
Despite these limitations, the service is functional for basic IP country filtering. It can be improved and extended as needed – for example, adding per-customer allowlists, integrating a UI for management, using an API gateway to handle external access, or improving the error responses with JSON output. For now, it provides a simple, containerized solution for country-based IP filtering.

**The auto-update process requires a license key, but I don't have access to it. So, I implemented the auto-update logic but left it disabled.**
