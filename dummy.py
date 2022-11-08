from locust import User, task


class Locust(User):
	@task
	def hello(self):
		pass
