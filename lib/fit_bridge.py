import sys
import json
import numpy as np
from scipy.optimize import curve_fit
 

def generalized_logistic_function(x, A, K, B, v, Q):
    return A + (K - A) / (1 + Q * np.exp(-B * (x - v)))

def fit_curve(x_data, y_data):
    try:
        x_data = np.array(x_data)
        y_data = np.array(y_data)
        
        # Sort data
        sort_idx = np.argsort(x_data)
        x_data, y_data = x_data[sort_idx], y_data[sort_idx]

        # Calculate A (Lower Asymptote)
        observed_min, observed_max = np.min(y_data), np.max(y_data)
        total_range = (observed_max - observed_min) * (4/3)
        A = float(round(observed_min - total_range * 0.25))

        # Initial guesses and bounds
        p0 = [np.max(y_data) * 1.05, 0.05, np.median(x_data), 1.0]

        # Adjusted bounds for classic totals (K max up to 1.5x max total)
        bounds = (
            [np.max(y_data) * 0.8, 0.001, np.min(x_data), 0.01],
            [np.max(y_data) * 1.5, 2.0,   np.max(x_data), 100.0]
        )

        popt, _ = curve_fit(
            lambda x, K, B, v, Q: generalized_logistic_function(x, A, K, B, v, Q),
            x_data, y_data, p0=p0, bounds=bounds, loss="soft_l1", maxfev=20000, ftol=1e-5,
        )
        
        return {"A": A, "K": popt[0], "B": popt[1], "V": popt[2], "Q": popt[3]}
    except Exception as e:
        return {"error": str(e)}

if __name__ == "__main__":
    # Read data from Go via Stdin
    input_data = json.load(sys.stdin)
    result = fit_curve(input_data['x'], input_data['y'])
    print(json.dumps(result))