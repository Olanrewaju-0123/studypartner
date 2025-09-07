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
        <p className="text-gray-500">No flashcards available</p>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto">
      {/* Progress */}
      <div className="mb-6">
        <div className="flex justify-between items-center mb-2">
          <span className="text-sm text-gray-600">
            Card {currentIndex + 1} of {flashcards.length}
          </span>
          <button
            onClick={resetCards}
            className="text-sm text-blue-600 hover:text-blue-700 flex items-center"
          >
            <RotateCcw className="h-4 w-4 mr-1" />
            Reset
          </button>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-2">
          <div
            className="bg-blue-600 h-2 rounded-full transition-all duration-300"
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
            bg-white rounded-xl shadow-lg p-8 min-h-[300px] cursor-pointer transition-transform duration-500 transform-gpu
            ${isFlipped ? "rotate-y-180" : ""}
          `}
          onClick={flipCard}
          style={{ transformStyle: "preserve-3d" }}
        >
          <div className={`${isFlipped ? "hidden" : "block"}`}>
            <div className="text-center">
              <h3 className="text-lg font-semibold text-gray-600 mb-4">
                Question
              </h3>
              <p className="text-xl text-gray-900 leading-relaxed">
                {currentCard.question}
              </p>
            </div>
          </div>

          <div className={`${isFlipped ? "block" : "hidden"}`}>
            <div className="text-center">
              <h3 className="text-lg font-semibold text-gray-600 mb-4">
                Answer
              </h3>
              <p className="text-xl text-gray-900 leading-relaxed">
                {currentCard.answer}
              </p>
            </div>
          </div>
        </div>

        <div className="text-center mt-4">
          <p className="text-sm text-gray-500">
            Click the card to {isFlipped ? "see question" : "reveal answer"}
          </p>
        </div>
      </div>

      {/* Navigation */}
      <div className="flex justify-between items-center mt-8">
        <button
          onClick={prevCard}
          disabled={currentIndex === 0}
          className={`
            flex items-center px-4 py-2 rounded-lg transition-colors
            ${
              currentIndex === 0
                ? "text-gray-400 cursor-not-allowed"
                : "text-blue-600 hover:bg-blue-50"
            }
          `}
        >
          <ChevronLeft className="h-5 w-5 mr-1" />
          Previous
        </button>

        <button
          onClick={nextCard}
          disabled={currentIndex === flashcards.length - 1}
          className={`
            flex items-center px-4 py-2 rounded-lg transition-colors
            ${
              currentIndex === flashcards.length - 1
                ? "text-gray-400 cursor-not-allowed"
                : "text-blue-600 hover:bg-blue-50"
            }
          `}
        >
          Next
          <ChevronRight className="h-5 w-5 ml-1" />
        </button>
      </div>
    </div>
  );
}
