package web

const loginPageHTML = `
<!DOCTYPE html>
<html>
<head><title>RxDiet Secure Portal - Login</title></head>
<body style="font-family: Arial, sans-serif; margin: 50px;">
    <h2>RxDiet Core Portal Login</h2>
    %s
    <form method="POST" action="/login">
        <label>Username:</label><br>
        <input type="text" name="username" required><br><br>
        <label>Password:</label><br>
        <input type="password" name="password" required><br><br>
        <button type="submit">Login</button>
    </form>
</body>
</html>
`

const dashboardHTML = `
<!DOCTYPE html>
<html>
<head><title>Dashboard</title></head>
<body style="font-family: Arial, sans-serif; margin: 50px;">
    <h2 style="color: green;">Welcome to the Secure Dashboard!</h2>
    <p>User Identity Status: <strong>Authenticated as %s</strong></p>
    <p>This data is protected behind server-side session tokens.</p>
    <a href="/logout">Logout</a>
</body>
</html>
`

const loginFailHTML = `
<!DOCTYPE html>
<html>
<head><title>403 - Unauthorized</title></head>
<body style="font-family: Arial, sans-serif; margin: 50px; text-align: center;">
    <h1 style="color: red;">403 - Access Denied</h1>
    <p>You must be logged in to view this page resource.</p>
    <p><a href="/login">Return to Login Page</a></p>
</body>
</html>
`