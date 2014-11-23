package com.ganddev.devfest2014.api;

import android.util.Log;

import com.android.volley.Request;
import com.android.volley.Response;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonArrayRequest;
import com.android.volley.toolbox.JsonObjectRequest;
import com.ganddev.devfest2014.DevFestApp;
import com.ganddev.devfest2014.R;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

/**
 * Created by bjornahlfeld on 22.11.14.
 */
public class ApiService {


    private static final String TAG = ApiService.class.getSimpleName();
    private static final DevFestApp ctx = DevFestApp.getInstance();

    public static void postGCMRegId(final String regId)
    {
        JSONObject obj = new JSONObject();
        try {
            obj.put("devicetoken", regId);
        } catch (JSONException e) {
            Log.e(TAG, e.getMessage());
        }

        JsonObjectRequest request = new JsonObjectRequest(Request.Method.POST,ctx.getString(R.string.api_gcm),obj, new Response.Listener<JSONObject>() {
            @Override
            public void onResponse(JSONObject response) {
                if(response != null) {
                    Log.i(TAG, "Response: " + response.toString());
                }
            }
        }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                if(error != null && error.getMessage() != null) {
                    Log.e(TAG, error.getMessage());
                }
            }
        });
        ctx.addToRequestQueue(request);
    }

    /**
     * Query the api for all articles...
     */
    public static void getArticles() {
        final JsonArrayRequest req = new JsonArrayRequest(ctx.getString(R.string.api_list_articles), new Response.Listener<JSONArray>() {
            @Override
            public void onResponse(JSONArray response) {

            }
        }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                if(error != null && error.getMessage() != null) {
                    Log.e(TAG, error.getMessage());
                }
            }
        });
        ctx.addToRequestQueue(req);
    }

    /**
     * Query the api for one article...
     * @param articleId
     */
    public static void getArticle(int articleId) {
        final JsonObjectRequest req = new JsonObjectRequest(ctx.getString(R.string.api_list_articles)+"/"+articleId, null, new Response.Listener<JSONObject>() {
            @Override
            public void onResponse(JSONObject response) {

            }
        }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                if(error != null && error.getMessage() != null)
                {
                    Log.e(TAG, error.getMessage());
                }
            }
        });
        ctx.addToRequestQueue(req);
    }
}
