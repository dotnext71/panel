import bottle


@bottle.route("/")
def index():
    return "Hello from Bottle on DotNext."


# gunicorn serves the module-level WSGI application (app:app).
app = bottle.default_app()
