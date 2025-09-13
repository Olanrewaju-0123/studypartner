"use client";

import { useState } from "react";
import {
  CheckCircle,
  XCircle,
  RotateCcw,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { Quiz as QuizType } from "@/types";

interface QuizProps {
  quiz: QuizType[];
}

export default function Quiz({ quiz }: QuizProps) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [selectedAnswers, setSelectedAnswers] = useState<number[]>(
    new Array(quiz.length).fill(-1)
  );
  const [showResults, setShowResults] = useState(false);
  const [quizCompleted, setQuizCompleted] = useState(false);

  const currentQuestion = quiz[currentIndex];

  // Debug logging
  console.log("Quiz component received:", quiz);
  console.log("Current question:", currentQuestion);
  if (currentQuestion) {
    console.log("Current question options:", currentQuestion.options);
  }

  const handleAnswerSelect = (answerIndex: number) => {
    if (showResults || quizCompleted) return;

    const newAnswers = [...selectedAnswers];
    newAnswers[currentIndex] = answerIndex;
    setSelectedAnswers(newAnswers);
  };

  const nextQuestion = () => {
    if (currentIndex < quiz.length - 1) {
      setCurrentIndex(currentIndex + 1);
      setShowResults(false); // Reset showResults when moving to next question
    } else {
      setQuizCompleted(true);
    }
  };

  const prevQuestion = () => {
    if (currentIndex > 0) {
      setCurrentIndex(currentIndex - 1);
      setShowResults(false); // Reset showResults when moving to previous question
    }
  };

  const showAnswer = () => {
    setShowResults(true);
  };

  const resetQuiz = () => {
    setCurrentIndex(0);
    setSelectedAnswers(new Array(quiz.length).fill(-1));
    setShowResults(false);
    setQuizCompleted(false);
  };

  const calculateScore = () => {
    let correct = 0;
    selectedAnswers.forEach((answer, index) => {
      if (answer === quiz[index].answer) {
        correct++;
      }
    });
    return Math.round((correct / quiz.length) * 100);
  };

  if (!currentQuestion) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-500">No quiz questions available</p>
      </div>
    );
  }

  if (quizCompleted) {
    const score = calculateScore();
    return (
      <div className="max-w-2xl mx-auto text-center">
        <div className="bg-white rounded-xl shadow-lg p-8">
          <div className="mb-6">
            <div
              className={`w-20 h-20 mx-auto rounded-full flex items-center justify-center mb-4 ${
                score >= 70
                  ? "bg-green-100"
                  : score >= 50
                  ? "bg-yellow-100"
                  : "bg-red-100"
              }`}
            >
              <span
                className={`text-2xl font-bold ${
                  score >= 70
                    ? "text-green-600"
                    : score >= 50
                    ? "text-yellow-600"
                    : "text-red-600"
                }`}
              >
                {score}%
              </span>
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-2">
              Quiz Complete!
            </h2>
            <p className="text-gray-600">
              You scored {score}% (
              {
                selectedAnswers.filter(
                  (answer, index) => answer === quiz[index].answer
                ).length
              }{" "}
              out of {quiz.length} correct)
            </p>
          </div>

          <button
            onClick={resetQuiz}
            className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors flex items-center mx-auto"
          >
            <RotateCcw className="h-5 w-5 mr-2" />
            Retake Quiz
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto">
      {/* Progress */}
      <div className="mb-6">
        <div className="flex justify-between items-center mb-2">
          <span className="text-sm text-gray-600">
            Question {currentIndex + 1} of {quiz.length}
          </span>
          <button
            onClick={resetQuiz}
            className="text-sm text-blue-600 hover:text-blue-700 flex items-center"
          >
            <RotateCcw className="h-4 w-4 mr-1" />
            Reset
          </button>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-2">
          <div
            className="bg-blue-600 h-2 rounded-full transition-all duration-300"
            style={{ width: `${((currentIndex + 1) / quiz.length) * 100}%` }}
          />
        </div>
      </div>

      {/* Question */}
      <div className="bg-white rounded-xl shadow-lg p-8 mb-6">
        <h3 className="text-xl font-semibold text-gray-900 mb-6">
          {currentQuestion.question}
        </h3>

        <div className="space-y-3">
          {currentQuestion.options && currentQuestion.options.length > 0 ? (
            currentQuestion.options.map((option, index) => {
              const isSelected = selectedAnswers[currentIndex] === index;
              const isCorrect = index === currentQuestion.answer;
              const isWrong = isSelected && !isCorrect;

              return (
                <button
                  key={index}
                  onClick={() => handleAnswerSelect(index)}
                  disabled={showResults}
                  className={`
                    w-full text-left p-4 rounded-lg border-2 transition-all
                    ${
                      showResults
                        ? isCorrect
                          ? "border-green-500 bg-green-50 text-green-900"
                          : isWrong
                          ? "border-red-500 bg-red-50 text-red-900"
                          : "border-gray-200 bg-gray-50 text-gray-600"
                        : isSelected
                        ? "border-blue-500 bg-blue-50 text-blue-900"
                        : "border-gray-200 hover:border-gray-300 hover:bg-gray-50"
                    }
                    ${showResults ? "cursor-default" : "cursor-pointer"}
                  `}
                >
                  <div className="flex items-center justify-between">
                    <span className="font-medium text-gray-900">
                      {String.fromCharCode(65 + index)}. {option || "No option text"}
                    </span>
                    {showResults && isCorrect && (
                      <CheckCircle className="h-5 w-5 text-green-500" />
                    )}
                    {showResults && isWrong && (
                      <XCircle className="h-5 w-5 text-red-500" />
                    )}
                  </div>
                </button>
              );
            })
          ) : (
            <div className="text-center py-8">
              <p className="text-red-500">No quiz options available for this question.</p>
              <p className="text-sm text-gray-500 mt-2">
                This might be a data issue. Please try regenerating the quiz.
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Navigation */}
      <div className="flex justify-between items-center">
        <button
          onClick={prevQuestion}
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

        {!showResults && selectedAnswers[currentIndex] !== -1 && (
          <button
            onClick={showAnswer}
            className="bg-green-600 text-white px-6 py-2 rounded-lg hover:bg-green-700 transition-colors"
          >
            Check Answer
          </button>
        )}

        <button
          onClick={nextQuestion}
          disabled={selectedAnswers[currentIndex] === -1}
          className={`
            flex items-center px-4 py-2 rounded-lg transition-colors
            ${
              selectedAnswers[currentIndex] === -1
                ? "text-gray-400 cursor-not-allowed"
                : "text-blue-600 hover:bg-blue-50"
            }
          `}
        >
          {currentIndex === quiz.length - 1 ? "Finish" : "Next"}
          <ChevronRight className="h-5 w-5 ml-1" />
        </button>
      </div>
    </div>
  );
}
