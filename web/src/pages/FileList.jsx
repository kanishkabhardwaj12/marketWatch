import React, { useState } from 'react';

function FileList() {
  const [isVisible, setIsVisible] = useState(false);

  const [fetchState, setFetchState] = useState({
    status: 'idle', 
    files: [],
    error: null,
  });

  const handleButtonClick = async () => {
    if (isVisible && fetchState.status === 'success') {
      setIsVisible(false);
      return;
    }
    
    setIsVisible(true);

    if (fetchState.status === 'idle' || fetchState.status === 'error') {
      setFetchState({ status: 'pending', files: [], error: null });
      try {
        const response = await fetch('/api/listFiles');
        if (!response.ok) {
          throw new Error('Failed to fetch files from the server.');
        }
        const data = await response.json();
        setFetchState({ status: 'success', files: data || [], error: null });
      } catch (err) {
        setFetchState({ status: 'error', files: [], error: err });
      }
    }
  };

  const getButtonText = () => {
    if (fetchState.status === 'pending') return 'Loading...';
    if (fetchState.status === 'error') return 'Retry';
    if (isVisible) return 'Hide Files';
    return 'Show Uploaded Files';
  };

  let content;
  if (fetchState.status === 'pending') {
    content = <p>Loading files</p>;
  } else if (fetchState.status === 'error') {
    content = <p className="error-message">Error: {fetchState.error.message}</p>;
  } else if (fetchState.status === 'success') {
    content = fetchState.files.length === 0 ? (
      <p>No files uploaded yet.</p>
    ) : (
      <>
        <h3>Uploaded Files:</h3>
        <ul>
          {fetchState.files.map((file) => (
            <li key={file.file_path}>
              <a href={`/tradebook_uploads/${file.file_path}`} target="_blank" rel="noopener noreferrer">
                {file.filename}
              </a>
            </li>
          ))}
        </ul>
      </>
    );
  }

  return (
    <div className="section">
      <button 
        type="button" 
        className="btn" 
        onClick={handleButtonClick}
        disabled={fetchState.status === 'pending'}
      >
        {getButtonText()}
      </button>

      {isVisible && (
        <div id="fileList" className="section">
          {content}
        </div>
      )}
    </div>
  );
}

export default FileList;