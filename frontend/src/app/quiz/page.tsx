"use client";

import React, { useActionState } from "react";
import { AllWords, TranslationKey } from "../types/word";
import { Button } from "@mui/material";
import { QuizButtonList, QuizContainer } from "./page.style";
import { OverridableStringUnion } from "@mui/types";
import { selectAnswer } from "../lib/actions";

const t_strings: TranslationKey[] = ['t_1', 't_2', 't_3', 't_4']

export default function Quiz() {
    const [allWords, setAllWords] = React.useState<AllWords>();

    const [currentWord, setCurrentWord] = React.useState(0);
    const [wasCorrect, setWasCorrect] = React.useState(false);
    const [score, setScore] = React.useState(0);

    React.useEffect(() => {
        fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/all-words`)
            .then((res) => res.json())
            .then((data) => setAllWords(data));
    }, []);

    const handleSelect = async (selectedWord: string, wordId: number) => {
        selectAnswer(wordId, selectedWord);
        if (selectedWord === allWords?.allWords[currentWord]?.t_1) {
            setScore(score + 1);
            setWasCorrect(true);
            setTimeout(() => {
                setCurrentWord(currentWord + 1);
                setWasCorrect(false);
            }, 1000);
        }
    };

    function shuffleArray(array: TranslationKey[]) {
        for (let i = array.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1)); // Random index from 0 to i
            [array[i], array[j]] = [array[j], array[i]]; // Swap elements
        }
        return array;
    }

    const [indices, setIndices] = React.useState(shuffleArray(t_strings))

    React.useEffect(() => {
        setIndices(shuffleArray(t_strings))
    }, [currentWord]);

    return (
        <QuizContainer>
            {allWords && allWords.allWords && allWords?.allWords.length > 0 && (
                <>
                    <h2 style={{ marginBottom: "50px" }}>
                        Current word: {allWords?.allWords[currentWord || 0]?.word}
                        {wasCorrect && <span> - Correct!</span>}
                    </h2>
                    <QuizButtonList>
                    {indices.map((index: TranslationKey ) => (
                            <Button variant='contained'
                                color={wasCorrect && index === 't_1' ? 'success' : 'primary' as OverridableStringUnion<'success' | 'primary'>}
                                key={index}
                                onClick={() => handleSelect(allWords.allWords[currentWord || 0]?.[index], allWords.allWords[currentWord].id)}>
                                {allWords.allWords[currentWord || 0]?.[index]}
                            </Button>
                        ))}
                    </QuizButtonList>
                </>
            )}
        </QuizContainer>
    );
}
