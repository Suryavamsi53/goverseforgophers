with open('ui/templates/pages/practice_questions.html', 'r') as f:
    content = f.read()

content = content.replace(
    r"""<span x-text="revealed ? \\'Hide Solution\\' : \\'View Solution\\'"></span>""",
    r"""<span x-text="revealed ? \'Hide Solution\' : \'View Solution\'"></span>"""
)

with open('ui/templates/pages/practice_questions.html', 'w') as f:
    f.write(content)

print("Fixed line 155")
