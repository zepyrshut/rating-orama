from flask import Blueprint, jsonify

from logic.utils import get_tv_show_episodes

tv_show_bp = Blueprint('tv_show_bp', __name__)


@tv_show_bp.route('/tv-show/<tt_id>', methods=['GET'])
def get_tv_show(tt_id):
    try:
        tv_show = get_tv_show_episodes(tt_id)
        return jsonify(tv_show.__dict__)
    except Exception as e:
        return jsonify({'error': str(e)}), 500
