from contextlib import closing
import json
import base64
import os
import sqlite3
import requests
import pytest
import tempfile
from dotenv import load_dotenv
from psycopg2 import connect # postgres database adapter


BASE_URL = 'http://localhost:8081'
USERNAME = 'simulator'
PWD = 'super_safe!'
CREDENTIALS = ':'.join([USERNAME, PWD]).encode('ascii')
ENCODED_CREDENTIALS = base64.b64encode(CREDENTIALS).decode()
HEADERS = {'Connection': 'close',
           'Content-Type': 'application/json',
           f'Authorization': f'Basic {ENCODED_CREDENTIALS}'}


@pytest.fixture(scope="session")
def postgresql():
    # if ../.env exists, load it
    load_dotenv(override=False)
    dbHost = os.environ.get("DB_HOST")
    dbPort = os.environ.get("DB_PORT")
    dbUser = os.environ.get("DB_USER")
    dbPassword = os.environ.get("DB_PASS")
    dbName = os.environ.get("DB_TEST_NAME")

    conn = connect(
        dbname=dbName,
        user=dbUser,
        password=dbPassword,
        host=dbHost,
        port=dbPort
    )

    yield conn

    conn.close()

@pytest.fixture(scope="session")
def reset_db_rows(postgresql):
    cursor = postgresql.cursor()
    cursor.execute("DELETE FROM messages;")
    cursor.execute("DELETE FROM followers;")
    cursor.execute("DELETE FROM users;")
    cursor.execute("ALTER SEQUENCE messages_id_seq RESTART WITH 1;")
    cursor.execute("ALTER SEQUENCE users_id_seq RESTART WITH 1;")
    cursor.execute("ALTER SEQUENCE followers_id_seq RESTART WITH 1;")
    
    
    postgresql.commit()

def test_basic_no_db():
    url = f'{BASE_URL}/latest'
    response = requests.get(url)
    assert response.status_code == 200

def test_latest(reset_db_rows):
    # post something to updaet LATEST
    url = f"{BASE_URL}/register"
    data = {'username': 'test', 'email': 'test@test', 'pwd': 'foo'}
    params = {'latest': 1337}
    print("pre test_latest POST")
    response = requests.post(url, data=json.dumps(data),
                             params=params, headers=HEADERS)
    assert response.ok

    # verify that latest was updated
    url = f'{BASE_URL}/latest'
    response = requests.get(url, headers=HEADERS)
    assert response.ok
    assert response.json()['latest'] == 1337

def test_register(reset_db_rows):
    username = 'a'
    email = 'a@a.a'
    pwd = 'a'
    data = {'username': username, 'email': email, 'pwd': pwd}
    params = {'latest': 1}
    response = requests.post(f'{BASE_URL}/register',
                             data=json.dumps(data), headers=HEADERS, params=params)
    assert response.ok
    # TODO: add another assertion that it is really there

    # verify that latest was updated
    response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 1


def test_create_msg(reset_db_rows):
    username = 'a'
    data = {'content': 'Blub!'}
    url = f'{BASE_URL}/msgs/{username}'
    params = {'latest': 2}
    response = requests.post(url, data=json.dumps(data),
                             headers=HEADERS, params=params)
    assert response.ok

    # verify that latest was updated
    response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 2

def test_get_latest_user_msgs(reset_db_rows):
    username = 'a'

    query = {'no': 20, 'latest': 3}
    url = f'{BASE_URL}/msgs/{username}'
    response = requests.get(url, headers=HEADERS, params=query)
    assert response.status_code == 200

    got_it_earlier = False
    for msg in response.json():
        if msg['content'] == 'Blub!' and msg['user'] == username:
            got_it_earlier = True

    assert got_it_earlier

    # verify that latest was updated
    response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 3

def test_get_latest_msgs(reset_db_rows):
    username = 'a'
    query = {'no': 20, 'latest': 4}
    url = f'{BASE_URL}/msgs'
    response = requests.get(url, headers=HEADERS, params=query)
    assert response.status_code == 200

    got_it_earlier = False
    for msg in response.json():
        if msg['content'] == 'Blub!' and msg['user'] == username:
            got_it_earlier = True

    assert got_it_earlier

    # verify that latest was updated
    response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 4

def test_register_b(reset_db_rows):
    username = 'b'
    email = 'b@b.b'
    pwd = 'b'
    data = {'username': username, 'email': email, 'pwd': pwd}
    params = {'latest': 5}
    response = requests.post(f'{BASE_URL}/register', data=json.dumps(data),
                             headers=HEADERS, params=params)
    assert response.ok
    # TODO: add another assertion that it is really there

     # verify that latest was updated
    response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 5

def test_register_c(reset_db_rows):
    username = 'c'
    email = 'c@c.c'
    pwd = 'c'
    data = {'username': username, 'email': email, 'pwd': pwd}
    params = {'latest': 6}
    response = requests.post(f'{BASE_URL}/register', data=json.dumps(data),
                             headers=HEADERS, params=params)
    assert response.ok

    # verify that latest was updated
    response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 6

def test_follow_user(reset_db_rows):
    username = 'a'
    url = f'{BASE_URL}/fllws/{username}'
    data = {'follow': 'b'}
    params = {'latest': 7}
    response = requests.post(url, data=json.dumps(data),
                             headers=HEADERS, params=params)
    assert response.ok

    data = {'follow': 'c'}
    params = {'latest': 8}
    response = requests.post(url, data=json.dumps(data),
                             headers=HEADERS, params=params)
    assert response.ok

    query = {'no': 20, 'latest': 9}
    response = requests.get(url, headers=HEADERS, params=query)
    assert response.ok

    json_data = response.json()
    assert "b" in json_data["follows"]
    assert "c" in json_data["follows"]

    # verify that latest was updated
    response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 9

def test_a_unfollows_b(reset_db_rows):
    username = 'a'
    url = f'{BASE_URL}/fllws/{username}'

    #  first send unfollow command
    data = {'unfollow': 'b'}
    params = {'latest': 10}
    response = requests.post(url, data=json.dumps(data),
                             headers=HEADERS, params=params)
    assert response.ok

    # then verify that b is no longer in follows list
    query = {'no': 20, 'latest': 11}
    response = requests.get(url, params=query, headers=HEADERS)
    assert response.ok
    assert 'b' not in response.json()['follows']

    # verify that latest was updated
    response = requests.get(f'{BASE_URL}/latest', headers=HEADERS)
    assert response.json()['latest'] == 11
