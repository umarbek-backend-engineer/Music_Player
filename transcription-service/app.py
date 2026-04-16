from flask import Flask, request, jsonify
import whisper
import os
import uuid

app = Flask(__name__)

# 🔥 Better model (balance speed + accuracy)
model = whisper.load_model("large")
# use "medium" or "large" if you want higher accuracy (slower)

@app.route("/transcribe", methods=["POST"])
def transcribe():
    if "file" not in request.files:
        return jsonify({"error": "No file provided"}), 400

    file = request.files["file"]

    # ✅ unique temp file
    filename = f"{uuid.uuid4()}.mp3"
    filepath = os.path.join("/tmp", filename)

    file.save(filepath)

    try:
        # 🔥 force English (important for songs)
        result = model.transcribe(
            filepath,
            language="en",
            fp16=False  # avoids warning on CPU
        )

        lyrics = []

        for seg in result["segments"]:
            lyrics.append({
                "start": round(seg["start"], 2),
                "end": round(seg["end"], 2),
                "text": seg["text"].strip()
            })

        response = {
            "lyrics": lyrics,
            "language": result.get("language", "en"),
            "duration": round(result.get("duration", 0), 2)
        }

        return jsonify(response)

    except Exception as e:
        return jsonify({
            "error": "transcription failed",
            "details": str(e)
        }), 500

    finally:
        # ✅ cleanup
        if os.path.exists(filepath):
            os.remove(filepath)


# ✅ health check (for Docker)
@app.route("/health", methods=["GET"])
def health():
    return jsonify({"status": "ok"}), 200


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5001)