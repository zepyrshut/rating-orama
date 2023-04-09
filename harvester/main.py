from flask import Flask

from routes.ping import ping
from routes.tv_show_routes import tv_show_bp

app = Flask(__name__)

app.register_blueprint(tv_show_bp)
app.register_blueprint(ping)

version = '0.1.0'
appName = 'Rating Orama Harvester'
author = 'Pedro Pérez'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
