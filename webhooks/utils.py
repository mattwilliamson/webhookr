import json


def scalar_or_list(potential_list):
    if type(potential_list) is list:
        list_len = len(potential_list)

        if list_len == 0:
            return None

        if list_len == 1:
            return potential_list[0]

    return potential_list


def nice_dict(request_item):
    new_dict = {k: scalar_or_list(v) for k, v in request_item.iteritems()}
    return json.dumps(new_dict or None, indent=4, separators=(',', ': '))
