"use client";

import React from "react";
import { Word } from "../types/word";

export default function Quiz() {
  const [allWords, setAllWords] = React.useState<Word[]>([]);

  const [currentWord, setCurrentWord] = React.useState(0);
  const [wasCorrect, setWasCorrect] = React.useState(false);
  const [score, setScore] = React.useState(0);

  const useEffect = React.useEffect(() => {
    fetch("http://localhost:8080/all-words")
      .then((res) => res.json())
      .then((data) => setAllWords(data));
  }, []);

  const handleSelect = (selectedWord: string) => {
    if (selectedWord === allWords[currentWord].t_1) {
      setScore(score + 1);
      setWasCorrect(true);
      setTimeout(() => {
        setCurrentWord(currentWord + 1);
        setWasCorrect(false);
      }, 1000);
    }
  };

  return (
    <div>
      <h1 style={{ marginBottom: "50px" }}>Quiz!!</h1>
      {allWords.length > 0 && (
        <>
          <h2 style={{ marginBottom: "50px" }}>
            Current word: {allWords[currentWord || 0]?.word}
            {wasCorrect && <span> - Correct!</span>}
          </h2>
          <ul>
            <li onClick={() => handleSelect(allWords[currentWord || 0]?.t_1)}>
              {allWords[currentWord || 0]?.t_1}
            </li>
            <li onClick={() => handleSelect(allWords[currentWord || 0]?.t_2)}>
              {allWords[currentWord || 0]?.t_2}
            </li>
            <li onClick={() => handleSelect(allWords[currentWord || 0]?.t_3)}>
              {allWords[currentWord || 0]?.t_3}
            </li>
            <li onClick={() => handleSelect(allWords[currentWord || 0]?.t_4)}>
              {allWords[currentWord || 0]?.t_4}
            </li>
          </ul>
        </>
      )}
    </div>
  );
}
