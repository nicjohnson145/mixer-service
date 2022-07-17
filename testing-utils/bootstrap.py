import requests

HOST = 'http://localhost:30000'


def ok(resp):
    if resp.status_code != 200:
        print(resp.url)
        print(resp.text)
        exit(1)


FOO_USER = {'username': 'foo', 'password': 'foopass'}
BAR_USER = {'username': 'bar', 'password': 'barpass'}
BAZ_USER = {'username': 'baz', 'password': 'bazpass'}

# Register foo user
ok(requests.post(f"{HOST}/api/v1/auth/register-user", json=FOO_USER))
# Register bar user
ok(requests.post(f"{HOST}/api/v1/auth/register-user", json=BAR_USER))
# Register baz user
ok(requests.post(f"{HOST}/api/v1/auth/register-user", json=BAZ_USER))

# Login & give foo user 2 drinks
resp = requests.post(f"{HOST}/api/v1/auth/login", json=FOO_USER)
foo_token = resp.json()['access_token']
ok(requests.post(
    f"{HOST}/api/v1/drinks/create",
    headers={'MixerAuth': foo_token},
    json={
        'name': 'Daquiri',
        'primary_alcohol': 'Rum',
        'ingredients': [
            '2.5 oz white rum',
            '0.5 oz simple syrup',
            '1 oz lime',
        ],
        'publicity': 'public',
    },
))
ok(requests.post(
    f"{HOST}/api/v1/drinks/create",
    headers={'MixerAuth': foo_token},
    json={
        'name': 'Jack & Coke',
        'primary_alcohol': 'Whiskey',
        'ingredients': [
            '1 part Jack Daniels',
            '1 part Coke',
        ],
        'publicity': 'public',
    },
))

# Login & give bar user 1 drink
resp = requests.post(f"{HOST}/api/v1/auth/login", json=BAR_USER)
bar_token = resp.json()['access_token']
ok(requests.post(
    f"{HOST}/api/v1/drinks/create",
    headers={'MixerAuth': bar_token},
    json={
        'name': "Bee's Knees",
        'primary_alcohol': 'Gin',
        'ingredients': [
            '0.75 oz honey syrup',
            '0.75 oz lemon juice',
            '2 oz gin',
        ],
        'publicity': 'public',
    },
))

# Login & give bar user 1 drink
resp = requests.post(f"{HOST}/api/v1/auth/login", json=BAZ_USER)
baz_token = resp.json()['access_token']
ok(requests.post(
    f"{HOST}/api/v1/drinks/create",
    headers={'MixerAuth': baz_token},
    json={
        'name': "Old Fashioned",
        'primary_alcohol': 'Bourbon',
        'ingredients': [
            '2.5 oz bourbon',
            '2 dashes angostura bitters',
            '1 barspoon demerara syrup'
        ],
        'publicity': 'public',
    },
))
