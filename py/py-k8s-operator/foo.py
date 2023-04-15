

html = """
<html>
{foo}
</html>
"""

def lambda_func():
    h = html.format(html, foo="hello")
    print(h)


lambda_func()
