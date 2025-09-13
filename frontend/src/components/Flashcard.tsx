"use client";

import { useState } from "react";
import { RotateCcw, ChevronLeft, ChevronRight } from "lucide-react";
import { Flashcard as FlashcardType } from "@/types";

interface FlashcardProps {
  flashcards: FlashcardType[];
}

export default function Flashcard({ flashcards }: FlashcardProps) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isFlipped, setIsFlipped] = useState(false);

  const currentCard = flashcards[currentIndex];

  const nextCard = () => {
    if (currentIndex < flashcards.length - 1) {
      setCurrentIndex(currentIndex + 1);
      setIsFlipped(false);
    }
  };

  const prevCard = () => {
    if (currentIndex > 0) {
      setCurrentIndex(currentIndex - 1);
      setIsFlipped(false);
    }
  };

  const flipCard = () => {
    setIsFlipped(!isFlipped);
  };

  const resetCards = () => {
    setCurrentIndex(0);
    setIsFlipped(false);
  };

  if (!currentCard) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-500 dark:text-gray-400">No flashcards available</p>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto px-4 sm:px-0">
      {/* Progress */}
      <div className="mb-4 sm:mb-6">
        <div className="flex justify-between items-center mb-2">
          <span className="text-xs sm:text-sm text-gray-600 dark:text-gray-300">
            Card {currentIndex + 1} of {flashcards.length}
          </span>
          <button
            onClick={resetCards}
            className="text-xs sm:text-sm text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 flex items-center"
          >
            <RotateCcw className="h-3 w-3 sm:h-4 sm:w-4 mr-1" />
            Reset
          </button>
        </div>
        <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
          <div
            className="bg-blue-600 dark:bg-blue-500 h-2 rounded-full transition-all duration-300"
            style={{
              width: `${((currentIndex + 1) / flashcards.length) * 100}%`,
            }}
          />
        </div>
      </div>

      {/* Flashcard */}
      <div className="relative">
        <div
          className={`
            bg-white dark:bg-gray-800 rounded-xl shadow-lg p-6 sm:p-8 min-h-[250px] sm:min-h-[300px] cursor-pointer transition-transform duration-500 transform-gpu
            ${isFlipped ? "rotate-y-180" : ""}
          `}
          onClick={flipCard}
          style={{ transformStyle: "preserve-3d" }}
        >
          <div className={`${isFlipped ? "hidden" : "block"}`}>
            <div className="text-center">
              <h3 className="text-base sm:text-lg font-semibold text-gray-600 dark:text-gray-400 mb-3 sm:mb-4">
                Question
              </h3>
              <p className="text-lg sm:text-xl text-gray-900 dark:text-white leading-relaxed">
                {currentCard.question}
              </p>
            </div>
          </div>

          <div className={`${isFlipped ? "block" : "hidden"}`}>
            <div className="text-center">
              <h3 className="text-base sm:text-lg font-semibold text-gray-600 dark:text-gray-400 mb-3 sm:mb-4">
                Answer
              </h3>
              <p className="text-lg sm:text-xl text-gray-900 dark:text-white leading-relaxed">
                {currentCard.answer}
              </p>
            </div>
          </div>
        </div>

        <div className="text-center mt-3 sm:mt-4">
          <p className="text-xs sm:text-sm text-gray-500 dark:text-gray-400">
            Click the card to {isFlipped ? "see question" : "reveal answer"}
          </p>
        </div>
      </div>

      {/* Navigation */}
      <div className="flex justify-between items-center mt-6 sm:mt-8">
        <button
          onClick={prevCard}
          disabled={currentIndex === 0}
          className={`
            flex items-center px-3 sm:px-4 py-2 rounded-lg transition-colors text-sm sm:text-base
            ${
              currentIndex === 0
                ? "text-gray-400 dark:text-gray-500 cursor-not-allowed"
                : "text-blue-600 dark:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20"
            }
          `}
        >
          <ChevronLeft className="h-4 w-4 sm:h-5 sm:w-5 mr-1" />
          <span className="hidden sm:inline">Previous</span>
          <span className="sm:hidden">Prev</span>
        </button>

        <button
          onClick={nextCard}
          disabled={currentIndex === flashcards.length - 1}
          className={`
            flex items-center px-3 sm:px-4 py-2 rounded-lg transition-colors text-sm sm:text-base
            ${
              currentIndex === flashcards.length - 1
                ? "text-gray-400 dark:text-gray-500 cursor-not-allowed"
                : "text-blue-600 dark:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20"
            }
          `}
        >
          <span className="hidden sm:inline">Next</span>
          <span className="sm:hidden">Next</span>
          <ChevronRight className="h-4 w-4 sm:h-5 sm:w-5 ml-1" />
        </button>
      </div>
    </div>
  );
}
