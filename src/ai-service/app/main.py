"""
OneTasteFamily AI Service
"""
from fastapi import FastAPI

app = FastAPI(title="OneTasteFamily AI Service", version="1.0.0")


@app.get("/")
async def root():
    return {"message": "OneTasteFamily AI Service"}


@app.get("/health")
async def health():
    return {"status": "healthy"}

