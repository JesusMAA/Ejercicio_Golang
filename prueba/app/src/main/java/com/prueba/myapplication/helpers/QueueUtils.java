package com.prueba.myapplication.helpers;

import android.graphics.Bitmap;
import android.content.Context;
import android.util.LruCache;
import com.android.volley.Request;
import com.android.volley.RequestQueue;
import com.android.volley.toolbox.ImageLoader;
import com.android.volley.toolbox.Volley;

public class QueueUtils {
    private static QueueObject uniqueInstance;
    public static class QueueObject {
        private RequestQueue mRequestQueue;
        private ImageLoader mImageLoader;
        private static Context mCtx;
        private QueueObject(Context context) {
            mCtx = context;
            mRequestQueue = getRequestQueue();
            mImageLoader = new ImageLoader(mRequestQueue,
                    new ImageLoader.ImageCache() {
                        private final LruCache<String, Bitmap>
                                cache = new LruCache<String, Bitmap>(20);
                        @Override
                        public Bitmap getBitmap(String url) {
                            return cache.get(url);
                        }
                        @Override
                        public void putBitmap(String url, Bitmap bitmap) {
                            cache.put(url, bitmap);
                        }
                    });
        }

        public RequestQueue getRequestQueue() {
            if (mRequestQueue == null) {
                mRequestQueue = Volley.newRequestQueue(mCtx.getApplicationContext());
            }
            return mRequestQueue;
        }

        public <T> void addToRequestQueue(Request<T> req) {
            getRequestQueue().add(req);
        }

        public ImageLoader getImageLoader() {
            return mImageLoader;
        }
    }

    public static synchronized QueueObject getInstance(Context context) {
        if (uniqueInstance == null) {
            uniqueInstance = new QueueObject(context);
        }
        return uniqueInstance;
    }
}