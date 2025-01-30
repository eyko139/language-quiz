"use client";

import React, {useState} from "react";
import {Word} from "../types/word";
import Link from "next/link";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import Box from "@mui/material/Box";
import {Button, TextField} from "@mui/material";
import {MainContainer} from "./main.styles";
import LoadingSpinner from "@/app/components/partials/LoadingSpinnter";

export default function Home() {
    const [trans, setTrans] = useState<string>("");
    const [savedWords, setSavedWords] = useState<string[]>([]);

    const [allWords, setAllWords] = React.useState<{ allWords: Word[], percentageUntranslated: number, totalWords: number }>();

    const [translationLoading, setTranslationLoading] = useState<boolean>(false);
    const [translationSuccess, setTranslationSuccess] = useState<boolean>(false);

    React.useEffect(() => {
        fetch("http://10.212.46.13:8080/all-words")
            .then((res) => res.json())
            .then((data) => setAllWords(data));
    }, [savedWords]);

    const onSubmit = (e: any) => {
        e.preventDefault();
        setTrans("");
        fetch("http://10.212.46.13:8080/save-word", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({message: trans}),
        })
            .then((res) => res.text())
            .then((data) => setSavedWords([...savedWords, data]));
    };

    const getTranslations = async () => {
        setTranslationLoading(true);
        const res = await fetch("http://10.212.46.13:8080/translate").then((res) =>
            res.json(),
        );

        if (res?.length > 0) {
            setTranslationSuccess(true);
        } else {
            setTranslationSuccess(false);
        }
        setTranslationLoading(false);
        fetch("http://10.212.46.13:8080/all-words")
            .then((res) => res.json())
            .then((data) => setAllWords(data));
    };

    return (
        <MainContainer>
            <TextField onChange={input => setTrans(input.target.value)} id="outlined-basic" label="Outlined"
                       variant="outlined"/>
            <Button onClick={onSubmit} color='success' variant="contained">Add Word</Button>
            <div style={{marginBottom: "50px"}}>Words currently in database: {allWords?.totalWords}
            </div>
            <div style={{marginBottom: "50px"}}>Percentage Untranslated: {allWords?.percentageUntranslated}%
            </div>
            {allWords?.percentageUntranslated > 0 && (
                <Button onClick={getTranslations} color='success' variant="contained">Translate</Button>
            )}
            {translationLoading && (
                <div>
                    Translation in progress...
                    <LoadingSpinner/>
                </div>
            )}
            {translationSuccess && (
                <Link href="/quiz">
                    Translation successful! Click here to start Quiz
                </Link>
            )}
        </MainContainer>
    );
}
