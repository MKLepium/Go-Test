# I am fully aware that this is awfull, loading the model everytime I call it is not efficient at all.

import sys
import uuid
from diffusers import StableDiffusionPipeline, EulerDiscreteScheduler
import torch

def generate_image(prompt):
    # Generate a UUID for this image
    image_id = str(uuid.uuid4())

    model_id = "stabilityai/stable-diffusion-2"
    scheduler = EulerDiscreteScheduler.from_pretrained(model_id, subfolder="scheduler")
    pipe = StableDiffusionPipeline.from_pretrained(model_id, scheduler=scheduler, torch_dtype=torch.float16)
    pipe = pipe.to("cuda")

    image = pipe(prompt).images[0]
    # Save the image using the UUID as the filename
    file_path = f"{image_id}.png"
    image.save(file_path)
    print("I need something to find the output")
    print(image_id)
    print("I need something to find the output")
    return image_id

if __name__ == "__main__":
    prompt = " ".join(sys.argv[1:])  # Takes all arguments as a single prompt
    generate_image(prompt)
