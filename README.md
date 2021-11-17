# GWI assesment

## Solution description
 - The implementation's basic services are user service and asset service.
 - User service is responsible for creating, updating, deleting, etc a user entity.
 - Asset service is responsible for all relative crud operation per asset. It is also responsible for starring an asset per user request.
 - The storage solution consists of an in-memory concurrent map. So storage persistence is not supported.
 - User authentication is cookie based.

## Features
 - Create, Update, Get and Delete assets by using the path prefix: `/assets` and for a specific asset type e.x. charts, use `/assets/charts/chart/{id}`
 - A `GET` on a prefix either `assets` or an asset type e.x `/assets/charts/` lists all relevant data.
 - Use `POST` to create an asset entity, `PATCH` to update.
 - For example update insight with id 20 `PATCH`: endpoint `/assets/insights/insight{20}`.
 - A user can star a specific asset by using prefix `/starred-assets` followed by the asset type and the asset id, e.x. `/starred-assets//audience/{1}`.
 - Use `PUT` to star an asset and `DELETE` to unstar.
 - Use a `GET` request to list all user's favorites assets in `/starred-assets/`
 - All star related actions can be issued only be logged-in users.
 - User login and signup are supported in `/login` and `/login` respectively.
 - User related paths are under the prefix `/users/`. 
 - Create, Update, Get and Delete users in the following path `/users/user/{1}`.


## Install

### Get the app

- `git clone git@github.com:nmakro/platform2.0-go-challenge.git`


### Build the container
- To build and start:
- `make docker-build`
- `make docker-run`

- To stop the container from running:
- `make docker-down`

### Local build and execution
- To build and run:
- `make app`
- `make run`
- or for one combined command use `make serve`
- `make clen` to delete binary
- `make tests` runs the integration tests.
