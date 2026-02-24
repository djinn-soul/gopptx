//go:build ignore

// urlfetch_demo generates PPTX files from HTML strings using the urlfetch package.
//
// Run with: go run ./examples/34-urlfetch/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx/urlfetch"
)

const mlHTML = `<!DOCTYPE html>
<html>
<head>
  <title>Machine Learning Fundamentals</title>
  <meta name="description" content="An introduction to core concepts in machine learning and AI">
</head>
<body>
  <main>
    <h1>Machine Learning Fundamentals</h1>
    <p>Machine learning is a branch of artificial intelligence that enables systems to learn and improve from experience without being explicitly programmed.</p>

    <h2>Core Concepts</h2>
    <ul>
      <li>Supervised learning — training on labelled examples</li>
      <li>Unsupervised learning — finding hidden structure in data</li>
      <li>Reinforcement learning — learning through reward signals</li>
      <li>Feature engineering — transforming raw data into model inputs</li>
      <li>Overfitting and regularisation — keeping models generalisable</li>
    </ul>

    <h2>Common Algorithms</h2>
    <p>Practitioners choose algorithms based on problem type, dataset size, and interpretability requirements.</p>
    <ul>
      <li>Linear/Logistic Regression — fast, interpretable baselines</li>
      <li>Decision Trees and Random Forests — ensemble methods</li>
      <li>Support Vector Machines — high-dimensional classification</li>
      <li>Neural Networks — flexible function approximators</li>
      <li>k-Nearest Neighbours — instance-based learning</li>
    </ul>

    <h2>The ML Workflow</h2>
    <p>A typical machine learning project follows an iterative cycle from data collection through deployment.</p>
    <pre><code>1. Collect & clean data
2. Explore & visualise
3. Select & train model
4. Evaluate on held-out set
5. Tune hyperparameters
6. Deploy & monitor</code></pre>

    <h2>Evaluation Metrics</h2>
    <p>Choosing the right metric is as important as choosing the right algorithm. Common metrics include accuracy, precision, recall, F1-score, and AUC-ROC for classification; MAE, RMSE, and R² for regression.</p>

    <h2>Tools and Frameworks</h2>
    <ul>
      <li>Python — dominant language for ML research and production</li>
      <li>scikit-learn — classical ML algorithms and pipelines</li>
      <li>PyTorch / TensorFlow — deep learning frameworks</li>
      <li>Pandas / NumPy — data manipulation and numerical computing</li>
      <li>MLflow / Weights &amp; Biases — experiment tracking</li>
    </ul>

    <h2>Getting Started</h2>
    <p>The fastest way to begin is with a structured dataset and scikit-learn. Install the core scientific stack, then load a dataset and fit your first model.</p>
    <pre><code>pip install scikit-learn pandas matplotlib
from sklearn.datasets import load_iris
from sklearn.ensemble import RandomForestClassifier
X, y = load_iris(return_X_y=True)
clf = RandomForestClassifier().fit(X, y)
print(clf.score(X, y))</code></pre>
  </main>
</body>
</html>`

const apiHTML = `<!DOCTYPE html>
<html>
<head><title>API Documentation</title></head>
<body>
  <main>
    <h1>REST API Reference</h1>
    <p>This document describes the REST API endpoints available for integration with our platform.</p>

    <h2>Authentication</h2>
    <p>All API requests require authentication using Bearer tokens in the Authorization header.</p>
    <pre><code>Authorization: Bearer YOUR_API_KEY</code></pre>

    <h2>Endpoints</h2>
    <h3>GET /users</h3>
    <p>Retrieve a list of all users in the system with pagination support.</p>
    <ul>
      <li>page - Page number (default: 1)</li>
      <li>limit - Items per page (default: 20)</li>
      <li>sort - Sort field (name, email, created_at)</li>
    </ul>

    <h3>POST /users</h3>
    <p>Create a new user account with the specified details and permissions.</p>

    <h2>Error Codes</h2>
    <table>
      <tr><th>Code</th><th>Meaning</th></tr>
      <tr><td>200</td><td>Success</td></tr>
      <tr><td>400</td><td>Bad Request</td></tr>
      <tr><td>401</td><td>Unauthorized</td></tr>
      <tr><td>404</td><td>Not Found</td></tr>
      <tr><td>500</td><td>Server Error</td></tr>
    </table>

    <h2>Rate Limiting</h2>
    <p>API requests are limited to 100 requests per minute per API key to ensure fair usage.</p>
  </main>
</body>
</html>`

func main() {
	outDir := filepath.Join("examples", "output")
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
		os.Exit(1)
	}

	// Example 1: HTML string → PPTX with defaults.
	fmt.Println("📄 Example 1: HTML to PPTX (defaults)")
	pptx, err := urlfetch.HTMLToPPTX(mlHTML, "https://example.com/ml-fundamentals")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ❌ error: %v\n", err)
		os.Exit(1)
	}
	path := filepath.Join(outDir, "34_urlfetch_ml_intro.pptx")
	if err := os.WriteFile(path, pptx, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "  ❌ write: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✅ Created %s (%d bytes)\n\n", path, len(pptx))

	// Example 2: Custom config and options.
	fmt.Println("📄 Example 2: Custom config and options")
	cfg := urlfetch.DefaultConfig().
		WithMaxSlides(5).
		WithMaxBullets(4).
		WithCode(true)

	opts := urlfetch.DefaultConversionOptions().
		WithTitle("ML Quick Reference").
		WithAuthor("gopptx").
		WithSourceURL(true)

	pptx, err = urlfetch.HTMLToPPTXWithOptions(mlHTML, "https://example.com/ml-fundamentals", cfg, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ❌ error: %v\n", err)
		os.Exit(1)
	}
	path = filepath.Join(outDir, "34_urlfetch_ml_quick.pptx")
	if err := os.WriteFile(path, pptx, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "  ❌ write: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✅ Created %s (%d bytes)\n\n", path, len(pptx))

	// Example 3: Technical documentation with a real table slide.
	fmt.Println("📄 Example 3: Technical documentation (with table)")
	pptx, err = urlfetch.HTMLToPPTX(apiHTML, "https://api.example.com/docs")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ❌ error: %v\n", err)
		os.Exit(1)
	}
	path = filepath.Join(outDir, "34_urlfetch_api_docs.pptx")
	if err := os.WriteFile(path, pptx, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "  ❌ write: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✅ Created %s (%d bytes)\n\n", path, len(pptx))

	fmt.Println("=== Done ===")
	fmt.Println("Generated files in examples/output/:")
	fmt.Println("  - 34_urlfetch_ml_intro.pptx")
	fmt.Println("  - 34_urlfetch_ml_quick.pptx")
	fmt.Println("  - 34_urlfetch_api_docs.pptx")
}
