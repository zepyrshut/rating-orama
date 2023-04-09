from flask import jsonify, Blueprint

ping = Blueprint('ping', __name__)


@ping.route('/ping', methods=['GET'])
def ping_pong():
    return jsonify('pong!')
