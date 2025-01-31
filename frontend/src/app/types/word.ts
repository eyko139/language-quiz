export interface Word {
  id: number;
  word: string;
  t_1: TranslationKey;
  t_2: TranslationKey;
  t_3: TranslationKey;
  t_4: TranslationKey;
  time: Date;
}

export type TranslationKey = 't_1' | 't_2' | 't_3' | 't_4'

export interface AllWords { allWords: Word[], percentageUntranslated: number, totalWords: number };
