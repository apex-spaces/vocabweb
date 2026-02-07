'use client';

import { useState, useCallback } from 'react';
import { useDropzone } from 'react-dropzone';

interface VocabWord {
  word: string;
  definition: string;
  pos: string;
  cefr_level: string;
  context_sentence: string;
}

interface OCRResult {
  extracted_text: string;
  words: VocabWord[];
}

export default function OCRPage() {
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [result, setResult] = useState<OCRResult | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [selectedWords, setSelectedWords] = useState<Set<string>>(new Set());

  const handleImageUpload = async (file: File) => {
    setError(null);
    setResult(null);
    setSelectedWords(new Set());

    // Create preview
    const reader = new FileReader();
    reader.onload = (e) => {
      setImagePreview(e.target?.result as string);
    };
    reader.readAsDataURL(file);

    // Upload and analyze
    setIsAnalyzing(true);
    const formData = new FormData();
    formData.append('image', file);
    formData.append('level', 'A2'); // Default level

    try {
      const response = await fetch('/api/v1/ocr/analyze', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Analysis failed');
      }

      const data: OCRResult = await response.json();
      setResult(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setIsAnalyzing(false);
    }
  };

  const onDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles.length > 0) {
      handleImageUpload(acceptedFiles[0]);
    }
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'image/*': ['.png', '.jpg', '.jpeg', '.gif', '.webp']
    },
    maxFiles: 1,
    maxSize: 10 * 1024 * 1024, // 10MB
  });

  const toggleWordSelection = (word: string) => {
    const newSelected = new Set(selectedWords);
    if (newSelected.has(word)) {
      newSelected.delete(word);
    } else {
      newSelected.add(word);
    }
    setSelectedWords(newSelected);
  };

  const handleAddToWordbook = async () => {
    if (selectedWords.size === 0) return;
    
    // TODO: Implement API call to add words to wordbook
    alert(`Adding ${selectedWords.size} words to wordbook`);
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">📸 Image OCR Analysis</h1>

        {/* Upload Area */}
        <div
          {...getRootProps()}
          className={`border-2 border-dashed rounded-lg p-12 text-center cursor-pointer transition-colors ${
            isDragActive
              ? 'border-blue-500 bg-blue-500/10'
              : 'border-gray-600 hover:border-gray-500'
          }`}
        >
          <input {...getInputProps()} />
          <div className="text-6xl mb-4">📷</div>
          <p className="text-lg mb-2">
            {isDragActive ? 'Drop image here...' : 'Drag & drop an image, or click to select'}
          </p>
          <p className="text-sm text-gray-400">Supports: PNG, JPG, GIF, WebP (max 10MB)</p>
        </div>

        {/* Image Preview */}
        {imagePreview && (
          <div className="mt-8">
            <h2 className="text-xl font-semibold mb-4">Uploaded Image</h2>
            <img
              src={imagePreview}
              alt="Preview"
              className="max-w-full h-auto rounded-lg border border-gray-700"
            />
          </div>
        )}

        {/* Loading State */}
        {isAnalyzing && (
          <div className="mt-8 text-center">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-white"></div>
            <p className="mt-4 text-lg">Analyzing image...</p>
          </div>
        )}

        {/* Error Message */}
        {error && (
          <div className="mt-8 p-4 bg-red-500/20 border border-red-500 rounded-lg">
            <p className="text-red-300">❌ {error}</p>
          </div>
        )}

        {/* Results */}
        {result && (
          <div className="mt-8 space-y-6">
            {/* Extracted Text */}
            <div>
              <h2 className="text-xl font-semibold mb-4">📝 Extracted Text</h2>
              <div className="bg-gray-800 p-4 rounded-lg border border-gray-700">
                <p className="whitespace-pre-wrap">{result.extracted_text}</p>
              </div>
            </div>

            {/* Vocabulary Words */}
            <div>
              <h2 className="text-xl font-semibold mb-4">📚 Vocabulary Words</h2>
              {result.words.length === 0 ? (
                <p className="text-gray-400">No new vocabulary words found.</p>
              ) : (
                <div className="space-y-4">
                  {result.words.map((word, index) => (
                    <div
                      key={index}
                      onClick={() => toggleWordSelection(word.word)}
                      className={`p-4 rounded-lg border cursor-pointer transition-colors ${
                        selectedWords.has(word.word)
                          ? 'bg-blue-500/20 border-blue-500'
                          : 'bg-gray-800 border-gray-700 hover:border-gray-600'
                      }`}
                    >
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="flex items-center gap-3 mb-2">
                            <h3 className="text-lg font-semibold">{word.word}</h3>
                            <span className="px-2 py-1 text-xs bg-purple-500/20 text-purple-300 rounded">
                              {word.cefr_level}
                            </span>
                            <span className="text-sm text-gray-400">{word.pos}</span>
                          </div>
                          <p className="text-gray-300 mb-2">{word.definition}</p>
                          <p className="text-sm text-gray-400 italic">"{word.context_sentence}"</p>
                        </div>
                        <div className="ml-4">
                          {selectedWords.has(word.word) ? '✅' : '⬜'}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {/* Add to Wordbook Button */}
              {result.words.length > 0 && (
                <button
                  onClick={handleAddToWordbook}
                  disabled={selectedWords.size === 0}
                  className="mt-6 w-full py-3 px-6 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 disabled:cursor-not-allowed rounded-lg font-semibold transition-colors"
                >
                  Add {selectedWords.size > 0 ? `${selectedWords.size} ` : ''}Selected Words to Wordbook
                </button>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

