# Portfolio functions
Serverless functions compilation which are used by my own Portfolio

## Proposals
**Proposals package** targets to handle the **clients proposals** requested via contact form shown in my Portfolio. Currently, **Proposals package** exposes 1 function:

- *Register*: Takes the payload sent from contact form via *http* and proccess its content to save in `Notion` database
  - It was written in `Go v1.20.6`
  - It uses `Notion` as *database*
  - It was deployed as a serverless function in *Digital Ocean*

