

start-emqx:
	@docker run -d --name emqx -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083  emqx/emqx-enterprise:6.2.0

open-emqx-dashboard:
	@open http://localhost:18083/#/dashboard/overview