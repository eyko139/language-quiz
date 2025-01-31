import React from "react";
import { AllWords } from "@/app/types/word";
import AddWord from "./partials/AddWord";
import styles from './styles.module.css'


export default async function Home() {
    const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/all-words`);
    const initial: AllWords = await res.json();

    // The data will be available directly in the component
    const { allWords, percentageUntranslated, totalWords } = initial;

    // const getTranslations = async () => {
    //     const res = await fetch(`${process.env.BACKEND_URL}/translate`).then((res) =>
    //         res.json(),
    //     )

    //     if (res?.length > 0) {
    //         // setTranslationSuccess(true);
    //     } else {
    //         // setTranslationSuccess(false);
    //     }
    //     // setTranslationLoading(false);
    //     fetch(`${process.env.BACKEND_URL}/all-words`)
    //         .then((res) => res.json())
    //     // .then((data) => setAllWords(data));
    // };

    return (
        <div className={styles.container}
        >
            <AddWord />
            <div style={{ marginBottom: "50px" }}>Words currently in database: {totalWords}
            </div>
            <div style={{ marginBottom: "50px" }}>Percentage Untranslated: {percentageUntranslated}%
            </div>
        </div>
    );
}
