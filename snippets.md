# 3 Snippets to test best responses with a range of token limit, temperature, top-k and top-p values

```javascript
// Snippet 1 Balanced
balanced = {
    "max_tokens": 150,
    "temperature": 0.7,
    "top_p": 0.9,
    "presence_penalty": 0.1,
    "frequency_penalty": 0.1
}

response = await openai.ChatCompletion.acreate(
    model="gpt-4",
    messages=[
        {"role": "system", "content": formatted_prompt},
        {"role": "user", "content": query.text}
    ],
    **balanced_config
)

// Snippet 2  Concise
concise_config = {
    "max_tokens": 75,
    "temperature": 0.3,
    "top_p": 0.8,
    "presence_penalty": 0,
    "frequency_penalty": 0.2
}

response = await openai.ChatCompletion.acreate(
    model="gpt-4",
    messages=[
        {"role": "system", "content": formatted_prompt},
        {"role": "user", "content": query.text}
    ],
    **concise_config
)

// Snippet 3 Creative

creative_config = {
    "max_tokens": 250,
    "temperature": 1.0,
    "top_p": 1.0,
    "presence_penalty": 0.2,
    "frequency_penalty": 0.2
}

response = await openai.ChatCompletion.acreate(
    model="gpt-4",
    messages=[
        {"role": "system", "content": formatted_prompt},
        {"role": "user", "content": query.text}
    ],
    **creative_config
)
```


