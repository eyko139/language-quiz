"use server";

import { revalidatePath } from "next/cache";

export async function addWord(formData: FormData) {
  const word = formData.get("word");
  const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/save-word`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ message: word }),
  });
  // if (res.ok) {
  //   const transRes = await fetch(
  //     `${process.env.NEXT_PUBLIC_BACKEND_URL}/translate`,
  //   );
  //   console.log(transRes);
  // }
  revalidatePath("/home");
}

export async function selectAnswer(wordId: number, translation: string) {
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_BACKEND_URL}/submit-answer/${wordId}/${translation}`,
    {
      method: "POST",
    },
  );
  return res.text();
}
