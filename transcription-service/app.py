from flask import Flask, request
import whisper

app = Flask(__name__)
model = whisper.load_model("base")

@app.route("/transcribe", methods=["POST"])
def transcribe():
    file = request.files["file"]
    file.save("temp.mp3")

    result = model.transcribe("temp.mp3")

    lines = []
    for seg in result["segments"]:
        start = seg["start"]
        text = seg["text"]

        m = int(start // 60)
        s = start % 60

        lines.append(f"[{m:02d}:{s:05.2f}] {text}")

    return {"lyrics": "\n".join(lines)}

app.run(host="0.0.0.0", port=5001)