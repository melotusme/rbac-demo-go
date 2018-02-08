from flask import Flask
from flask_sqlalchemy import SQLAlchemy

import flask_admin as admin
from flask_admin.contrib import sqla


# Create application
app = Flask(__name__)


# Create dummy secrey key so we can use sessions
app.config['SECRET_KEY'] = '123456790'

# Create in-memory database
app.config['SQLALCHEMY_DATABASE_URI'] = 'mysql://root:@localhost/rbac'
app.config['SQLALCHEMY_ECHO'] = True
db = SQLAlchemy(app)

# Flask views
@app.route('/')
def index():
    return '<a href="/admin/">Click me to get to Admin!</a>'


class Role(db.Model):
    __tablename__ = 'roles'
    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    name = db.Column(db.String(50))

    def __str__(self):
        return self.desc


# class RoleUsers(db.Model):
#     __tablename__="user_roles"
#     user_id = db.Column(db.Integer())
#     role_id = db.Column(db.Integer())

# Define models
user_roles = db.Table(
    'user_roles',
    db.Column('user_id', db.Integer(), db.ForeignKey('users.id')),
    db.Column('role_id', db.Integer(), db.ForeignKey('roles.id'))
)

class User(db.Model):
    __tablename__ = 'users'
    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    name = db.Column(db.String(50))
    roles = db.relationship(Role, secondary="user_roles",
                            backref=db.backref('users', lazy='dynamic'))


    def __str__(self):
        return self.desc

class Permission(db.Model):
    __tablename__ = 'permissions'
    id = db.Column(db.Integer, primary_key=True, autoincrement=True)
    route = db.Column(db.String(50))

    def __str__(self):
        return self.desc

class RoleAdmin(sqla.ModelView):
    column_display_pk = True
    form_columns = ['id', 'name']

class UserAdmin(sqla.ModelView):
    column_display_pk = True
    column_hide_backrefs = True
    form_columns = ('id', 'name')

class PermissionAdmin(sqla.ModelView):
    column_display_pk = True
    form_columns = ['id', 'route']


# Create admin
admin = admin.Admin(app, name='Permission Manager', template_mode='bootstrap3')
admin.add_view(RoleAdmin(Role, db.session))
admin.add_view(UserAdmin(User, db.session))
admin.add_view(PermissionAdmin(Permission, db.session))

if __name__ == '__main__':

    # Create DB
    # db.create_all()

    # Start app
    app.run(debug=True)
