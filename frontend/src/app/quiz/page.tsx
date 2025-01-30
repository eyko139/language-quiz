"use client";

import React from "react";
import {Word} from "../types/word";
import {Button} from "@mui/material";
import {QuizButtonList, QuizContainer} from "./quizPage.style";
import {OverridableStringUnion} from "@mui/types";

export default function Quiz() {
    const [allWords, setAllWords] = React.useState<{ allWords: Word[], percentageUntranslated: number }>();

    const [currentWord, setCurrentWord] = React.useState(0);
    const [wasCorrect, setWasCorrect] = React.useState(false);
    const [score, setScore] = React.useState(0);

    const useEffect = React.useEffect(() => {
        fetch("http://10.212.46.13:8080/all-words")
            .then((res) => res.json())
            .then((data) => setAllWords(data));
    }, []);

    const handleSelect = (selectedWord: string) => {
        if (selectedWord === allWords?.allWords[currentWord]?.t_1) {
            setScore(score + 1);
            setWasCorrect(true);
            setTimeout(() => {
                setCurrentWord(currentWord + 1);
                setWasCorrect(false);
            }, 1000);
        }
    };

    function shuffleArray(array) {
        for (let i = array.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1)); // Random index from 0 to i
            [array[i], array[j]] = [array[j], array[i]]; // Swap elements
        }
        return array;
    }

    const [indices, setIndices] = React.useState(shuffleArray(['t_1', 't_2', 't_3', 't_4']))

    React.useEffect(() => {
        setIndices(shuffleArray(['t_1', 't_2', 't_3', 't_4']))
    }, [currentWord]);

    return (
        <QuizContainer>
            {allWords && allWords.allWords && allWords?.allWords.length > 0 && (
                <>
                    <h2 style={{marginBottom: "50px"}}>
                        Current word: {allWords?.allWords[currentWord || 0]?.word}
                        {wasCorrect && <span> - Correct!</span>}
                    </h2>
                    <QuizButtonList>
                        {indices.map((index) => (
                            <Button variant='contained'
                                    color={wasCorrect && index === 't_1' ? 'success' : 'primary' as OverridableStringUnion<'success' | 'primary'>}
                                    key={index}
                                    onClick={() => handleSelect(allWords.allWords[currentWord || 0]?.[index])}>
                                {allWords.allWords[currentWord || 0]?.[index]}
                            </Button>
                        ))}
                    </QuizButtonList>
                </>
            )}
        </QuizContainer>
    );
}
