"use client";

import { useState } from "react";
import { ArrowLeft, FileText, Target, Brain, CheckCircle } from "lucide-react";
import Link from "next/link";

export default function DemoPage() {
  const [activeDemo, setActiveDemo] = useState<string | null>(null);

  const demoSample = {
    title: "Introduction to Machine Learning",
    content: `Machine learning is a subset of artificial intelligence that focuses on algorithms that can learn from data. There are three main types of machine learning:

1. Supervised Learning: Uses labeled training data to learn a mapping from inputs to outputs. Examples include classification and regression.

2. Unsupervised Learning: Finds hidden patterns in data without labeled examples. Examples include clustering and dimensionality reduction.

3. Reinforcement Learning: Learns through interaction with an environment, receiving rewards or penalties for actions. Used in game playing and robotics.

Key concepts include features, training data, models, and evaluation metrics. The goal is to build models that generalize well to new, unseen data.`,
  };

  const demoSummary = `Machine learning is an AI subset using algorithms to learn from data. It has three main types: supervised learning (using labeled data for classification/regression), unsupervised learning (finding patterns without labels), and reinforcement learning (learning through environment interaction). Key concepts include features, training data, models, and evaluation metrics, with the goal of building models that generalize to new data.`;

  const demoFlashcards = [
    {
      question: "What are the three main types of machine learning?",
      answer:
        "Supervised learning, unsupervised learning, and reinforcement learning.",
    },
    {
      question: "What is supervised learning?",
      answer:
        "Uses labeled training data to learn a mapping from inputs to outputs, including classification and regression.",
    },
    {
      question: "What is the goal of machine learning?",
      answer: "To build models that generalize well to new, unseen data.",
    },
  ];

  const demoQuiz = [
    {
      question: "Which type of machine learning uses labeled training data?",
      options: [
        "Unsupervised Learning",
        "Supervised Learning",
        "Reinforcement Learning",
        "Deep Learning",
      ],
      answer: 1,
    },
    {
      question: "What does reinforcement learning learn through?",
      options: [
        "Labeled data",
        "Hidden patterns",
        "Environment interaction",
        "Neural networks",
      ],
      answer: 2,
    },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center">
              <Link
                href="/"
                className="flex items-center text-gray-600 hover:text-gray-900"
              >
                <ArrowLeft className="h-5 w-5 mr-2" />
                Back to Home
              </Link>
            </div>
            <h1 className="text-xl font-semibold text-gray-900">Demo</h1>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold text-gray-900 mb-4">
            Try AI Study Partner
          </h2>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            Experience how AI can transform your study materials. Upload your
            own files or try our demo content.
          </p>
        </div>

        {/* Demo Content */}
        <div className="bg-white rounded-lg shadow-sm p-8 mb-8">
          <h3 className="text-2xl font-bold text-gray-900 mb-4">
            Sample Content: {demoSample.title}
          </h3>
          <div className="bg-gray-50 rounded-lg p-6 mb-6">
            <h4 className="font-semibold text-gray-900 mb-3">
              Original Content:
            </h4>
            <p className="text-gray-700 leading-relaxed whitespace-pre-line">
              {demoSample.content}
            </p>
          </div>

          {/* Demo Features */}
          <div className="grid md:grid-cols-3 gap-6">
            {/* Summary Demo */}
            <div className="border rounded-lg p-6">
              <div className="flex items-center mb-4">
                <FileText className="h-6 w-6 text-blue-600 mr-2" />
                <h4 className="text-lg font-semibold text-gray-900">
                  AI Summary
                </h4>
              </div>
              <p className="text-gray-600 text-sm mb-4">
                Get a concise summary of key concepts
              </p>
              <button
                onClick={() =>
                  setActiveDemo(activeDemo === "summary" ? null : "summary")
                }
                className="w-full bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
              >
                {activeDemo === "summary" ? "Hide" : "Show"} Summary
              </button>
              {activeDemo === "summary" && (
                <div className="mt-4 p-4 bg-blue-50 rounded-lg">
                  <p className="text-gray-700 text-sm leading-relaxed">
                    {demoSummary}
                  </p>
                </div>
              )}
            </div>

            {/* Flashcards Demo */}
            <div className="border rounded-lg p-6">
              <div className="flex items-center mb-4">
                <Target className="h-6 w-6 text-green-600 mr-2" />
                <h4 className="text-lg font-semibold text-gray-900">
                  Flashcards
                </h4>
              </div>
              <p className="text-gray-600 text-sm mb-4">
                Interactive Q&A cards for memorization
              </p>
              <button
                onClick={() =>
                  setActiveDemo(
                    activeDemo === "flashcards" ? null : "flashcards"
                  )
                }
                className="w-full bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition-colors"
              >
                {activeDemo === "flashcards" ? "Hide" : "Show"} Flashcards
              </button>
              {activeDemo === "flashcards" && (
                <div className="mt-4 space-y-3">
                  {demoFlashcards.map((card, index) => (
                    <div key={index} className="p-3 bg-green-50 rounded-lg">
                      <p className="font-medium text-gray-900 text-sm mb-1">
                        Q: {card.question}
                      </p>
                      <p className="text-gray-700 text-sm">A: {card.answer}</p>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Quiz Demo */}
            <div className="border rounded-lg p-6">
              <div className="flex items-center mb-4">
                <Brain className="h-6 w-6 text-purple-600 mr-2" />
                <h4 className="text-lg font-semibold text-gray-900">Quiz</h4>
              </div>
              <p className="text-gray-600 text-sm mb-4">
                Multiple choice questions to test knowledge
              </p>
              <button
                onClick={() =>
                  setActiveDemo(activeDemo === "quiz" ? null : "quiz")
                }
                className="w-full bg-purple-600 text-white px-4 py-2 rounded-lg hover:bg-purple-700 transition-colors"
              >
                {activeDemo === "quiz" ? "Hide" : "Show"} Quiz
              </button>
              {activeDemo === "quiz" && (
                <div className="mt-4 space-y-4">
                  {demoQuiz.map((question, index) => (
                    <div key={index} className="p-4 bg-purple-50 rounded-lg">
                      <p className="font-medium text-gray-900 text-sm mb-3">
                        {index + 1}. {question.question}
                      </p>
                      <div className="space-y-2">
                        {question.options.map((option, optIndex) => (
                          <div
                            key={optIndex}
                            className={`p-2 rounded text-sm ${
                              optIndex === question.answer
                                ? "bg-green-100 text-green-800 border border-green-300"
                                : "bg-white text-gray-700 border border-gray-200"
                            }`}
                          >
                            {String.fromCharCode(65 + optIndex)}. {option}
                            {optIndex === question.answer && (
                              <CheckCircle className="inline h-4 w-4 ml-2 text-green-600" />
                            )}
                          </div>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Call to Action */}
        <div className="text-center bg-blue-50 rounded-lg p-8">
          <h3 className="text-2xl font-bold text-gray-900 mb-4">
            Ready to Transform Your Study Materials?
          </h3>
          <p className="text-gray-600 mb-6">
            Upload your own documents and experience the power of AI-driven
            study tools.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              href="/upload"
              className="bg-blue-600 text-white px-8 py-3 rounded-lg text-lg font-semibold hover:bg-blue-700 transition-colors"
            >
              Upload Your Notes
            </Link>
            <Link
              href="/"
              className="border border-blue-600 text-blue-600 px-8 py-3 rounded-lg text-lg font-semibold hover:bg-blue-50 transition-colors"
            >
              Learn More
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
